package handlers

import (
	"github.com/didil/nginx-lb-updater/services"
	"go.uber.org/zap"
)

type App struct {
	lbUpdater services.LBUpdater
	logger    *zap.Logger
}

func NewApp(lbUpdater services.LBUpdater, logger *zap.Logger) *App {
	return &App{
		lbUpdater: lbUpdater,
		logger:    logger,
	}
}
