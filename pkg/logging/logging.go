// Package logging to configure logging
package logging

import (
	"github.com/shahariaazam/teredix/pkg/config"
	"github.com/sirupsen/logrus"
)

// Logging holds logging configuration
type Logging struct {
	appConfig *config.AppConfig
}

// NewLogger builds a new logging system
func NewLogger(appConfig *config.AppConfig) *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)

	loggingCfg := appConfig.Logging
	switch loggingCfg.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	switch loggingCfg.Level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warning":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
