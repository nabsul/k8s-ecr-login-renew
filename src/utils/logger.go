package utils

import "go.uber.org/zap"

// GetLogger function return zap sugar logger instance
// It also returns error if any error occur while instaniate zap logger
func GetLogger() (*zap.SugaredLogger, error) {
	devCfg := zap.NewProductionConfig()
	devCfg.DisableCaller = true
	devCfg.DisableStacktrace = true

	logger, err := devCfg.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	return logger.Sugar(), err
}
