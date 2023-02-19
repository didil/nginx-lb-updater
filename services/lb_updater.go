package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/pkg/errors"
)

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type LBUpdater interface {
	UpdateStream(backendName string, port int, protocol string, servers []Server, proxyTimeoutSeconds, proxyConnectTimeoutSeconds int) error
}

type lbUpdater struct {
	proxyTimeoutSeconds        int
	proxyConnectTimeoutSeconds int
	confFolderPath             string
	// synchronize updates
	lock *sync.Mutex
	tmpl *template.Template
}

func NewLBUpdater() (LBUpdater, error) {
	proxyTimeoutSeconds, err := strconv.Atoi(os.Getenv("PROXY_TIMEOUT_SECONDS"))
	if err != nil || proxyTimeoutSeconds <= 0 {
		proxyTimeoutSeconds = 5
	}

	proxyConnectTimeoutSeconds, err := strconv.Atoi(os.Getenv("PROXY_CONNECT_TIMEOUT_SECONDS"))
	if err != nil || proxyConnectTimeoutSeconds <= 0 {
		proxyConnectTimeoutSeconds = 2
	}

	confFolderPath := os.Getenv("CONF_FOLDER_PATH")
	if confFolderPath == "" {
		confFolderPath = "/etc/nginx/streams.d"
	}

	lock := &sync.Mutex{}

	tmpl, err := template.New("conf").Parse(confTemplate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile template")
	}

	return &lbUpdater{
		proxyTimeoutSeconds:        proxyTimeoutSeconds,
		proxyConnectTimeoutSeconds: proxyConnectTimeoutSeconds,
		confFolderPath:             confFolderPath,
		lock:                       lock,
		tmpl:                       tmpl,
	}, nil
}

var confTemplate string = `# generated by nginx-lb-updater, do not edit as changes can be overwritten
stream {
	upstream {{.BackendName}} {
		{{- range .UpstreamServers}}
		server {{.Host}}:{{.Port}};
		{{- end}}
	}

	server {
		listen        {{.Port}} {{.Protocol}};
		proxy_pass    {{.BackendName}};
		proxy_timeout {{.ProxyTimeoutSeconds}}s;
		proxy_connect_timeout {{.ProxyConnectTimeoutSeconds}}s;
	}
}
`

type LBUpdate struct {
	BackendName                string
	Port                       int
	Protocol                   string
	UpstreamServers            []Server
	ProxyTimeoutSeconds        int
	ProxyConnectTimeoutSeconds int
}

func (u *lbUpdater) UpdateStream(backendName string, port int, protocol string, servers []Server, proxyTimeoutSeconds, proxyConnectTimeoutSeconds int) error {
	// avoid concurrent reloads
	u.lock.Lock()
	defer u.lock.Unlock()

	backendName = strings.TrimSpace(backendName)
	if backendName == "" {
		return fmt.Errorf("empty backend name")
	}

	configFilename := filepath.Base(fmt.Sprintf("%s.conf", backendName))

	filename := filepath.Join(u.confFolderPath, configFilename)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to open conf file: %s", filename)
	}
	defer f.Close()

	lbUpdate := &LBUpdate{
		BackendName:                backendName,
		Port:                       port,
		UpstreamServers:            servers,
		ProxyTimeoutSeconds:        proxyTimeoutSeconds,
		ProxyConnectTimeoutSeconds: proxyConnectTimeoutSeconds,
	}

	// only needed for udb
	if strings.ToLower(protocol) == "udp" {
		lbUpdate.Protocol = "udp"
	}

	if lbUpdate.ProxyTimeoutSeconds <= 0 {
		// use default value
		lbUpdate.ProxyTimeoutSeconds = u.proxyTimeoutSeconds
	}

	if lbUpdate.ProxyConnectTimeoutSeconds <= 0 {
		// use default value
		lbUpdate.ProxyConnectTimeoutSeconds = u.proxyConnectTimeoutSeconds
	}

	err = u.tmpl.Execute(f, lbUpdate)
	if err != nil {
		return errors.Wrapf(err, "failed to write conf file: %s", filename)
	}

	return nil
}
