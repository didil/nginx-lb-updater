package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didil/nginx-lb-updater/mock"
	"github.com/didil/nginx-lb-updater/server"
	"github.com/didil/nginx-lb-updater/server/handlers"
	"github.com/didil/nginx-lb-updater/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestUpdateLB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lbUpdater := mock.NewMockLBUpdater(ctrl)
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	app := handlers.NewApp(lbUpdater, logger)
	r := server.NewRouter(app)
	s := httptest.NewServer(r)
	defer s.Close()

	requestBody := &handlers.UpdateRequest{
		BackendName: "namespace-a_myservice",
		LBPort:      9000,
		LBProtocol:  "tcp",
		UpstreamServers: []services.Server{
			{
				Host: "192.168.101.2",
				Port: 5014,
			},
			{
				Host: "192.168.101.3",
				Port: 5014,
			},
		},
		ProxyTimeoutSeconds:        6,
		ProxyConnectTimeoutSeconds: 3,
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(requestBody)
	assert.NoError(t, err)

	lbUpdater.EXPECT().
		UpdateStream(requestBody.BackendName, requestBody.LBPort, requestBody.LBProtocol, requestBody.UpstreamServers, requestBody.ProxyTimeoutSeconds, requestBody.ProxyConnectTimeoutSeconds).
		Return(nil)

	req, err := http.NewRequest(http.MethodPost, s.URL+"/api/v1/lb", &buf)
	assert.NoError(t, err)

	cl := &http.Client{}
	resp, err := cl.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	jsonOK := &handlers.JSONOK{}
	err = json.NewDecoder(resp.Body).Decode(jsonOK)
	assert.NoError(t, err)
}

func TestDeleteLB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lbUpdater := mock.NewMockLBUpdater(ctrl)
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	app := handlers.NewApp(lbUpdater, logger)
	r := server.NewRouter(app)
	s := httptest.NewServer(r)
	defer s.Close()

	requestBody := &handlers.DeleteRequest{
		BackendName: "namespace-a_myservice",
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(requestBody)
	assert.NoError(t, err)

	lbUpdater.EXPECT().
		DeleteStream(requestBody.BackendName).
		Return(nil)

	req, err := http.NewRequest(http.MethodDelete, s.URL+"/api/v1/lb", &buf)
	assert.NoError(t, err)

	cl := &http.Client{}
	resp, err := cl.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	jsonOK := &handlers.JSONOK{}
	err = json.NewDecoder(resp.Body).Decode(jsonOK)
	assert.NoError(t, err)
}
