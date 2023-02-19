package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/didil/nginx-lb-updater/services"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UpdateRequest struct {
	BackendName                string            `json:"backendName"`
	LBPort                     int               `json:"lbPort"`
	LBProtocol                 string            `json:"lbProtocol"`
	UpstreamServers            []services.Server `json:"upstreamServers"`
	ProxyTimeoutSeconds        int               `json:"proxyTimeoutSeconds"`
	ProxyConnectTimeoutSeconds int               `json:"proxyConnectTimeoutSeconds"`
}

func (app *App) UpdateLB(w http.ResponseWriter, r *http.Request) {
	req := &UpdateRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		app.WriteJSONErr(w, http.StatusBadRequest, errors.Wrapf(err, "failed to decode request"))
		return
	}

	err = app.lbUpdater.UpdateStream(req.BackendName, req.LBPort, req.LBProtocol, req.UpstreamServers, req.ProxyTimeoutSeconds, req.ProxyConnectTimeoutSeconds)
	if err != nil {
		app.WriteJSONErr(w, http.StatusInternalServerError, errors.Wrapf(err, "failed to update stream"))
		return
	}

	app.WriteJSONResponse(w, http.StatusOK, JSONOK{})
}

type JSONErr struct {
	Err error `json:"err"`
}

type JSONOK struct {
}

func (app *App) WriteJSONErr(w http.ResponseWriter, statusCode int, err error) {
	app.WriteJSONResponse(w, statusCode, &JSONErr{Err: err})
}

func (app *App) WriteJSONResponse(w http.ResponseWriter, statusCode int, resp any) {
	w.WriteHeader(statusCode)
	writeErr := json.NewEncoder(w).Encode(resp)
	if writeErr != nil {
		app.logger.Error("json write err", zap.Error(writeErr))
	}
}
