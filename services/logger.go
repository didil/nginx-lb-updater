package services

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	if os.Getenv("APP_ENV") == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize logger")
	}
	return logger, nil
}
