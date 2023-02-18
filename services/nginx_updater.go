package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type Server struct {
	Host string
	Port int
}

type NginxUpdater interface {
	UpdateStream(backendName string, port int, servers []Server) error
}

type nginxUpdater struct {
	proxyTimeout        string
	proxyConnectTimeout string
	confFolderPath      string
	// synchronize updates
	lock *sync.Mutex
}

func NewNginxUpdater() *nginxUpdater {
	proxyTimeout := os.Getenv("PROXY_TIMEOUT")
	if proxyTimeout == "" {
		proxyTimeout = "5s"
	}

	proxyConnectTimeout := os.Getenv("PROXY_CONNECT_TIMEOUT")
	if proxyConnectTimeout == "" {
		proxyConnectTimeout = "2s"
	}

	confFolderPath := os.Getenv("CONF_FOLDER_PATH")
	if confFolderPath == "" {
		confFolderPath = "/etc/nginx/conf.d"
	}

	lock := &sync.Mutex{}

	return &nginxUpdater{
		proxyTimeout:        proxyTimeout,
		proxyConnectTimeout: proxyConnectTimeout,
		confFolderPath:      confFolderPath,
		lock:                lock,
	}
}

func (u *nginxUpdater) UpdateStream(backendName string, port int, servers []Server) error {
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

	// TODO: write using template

	return nil
}
