package utils

import (
	"os"

	"go.uber.org/zap"
)

const (
	yukaLogEnv = "LOG_LEVEL"
)

// GetLogger builds a logger with generalised config
func GetLogger() (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	debug := os.Getenv(yukaLogEnv)
	if debug != "" {
		logCfg := zap.NewDevelopmentConfig()
		logCfg.EncoderConfig.TimeKey = ""
		logger, err = logCfg.Build()
		logger.Debug("Debug logging enabled")
	} else {
		logCfg := zap.NewProductionConfig()
		logCfg.DisableStacktrace = true
		logCfg.EncoderConfig.TimeKey = ""
		logger, err = logCfg.Build()
	}

	return logger, err
}
