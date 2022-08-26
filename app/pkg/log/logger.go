package logger

import (
	"fmt"

	"github.com/kadirgonen/movie-api/app/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(config *config.Config) {
	logLevel, err := zapcore.ParseLevel(config.Logger.Level)
	if err != nil {
		panic(fmt.Sprintf("Unknown log level %v", logLevel))
	}
	var cfg zap.Config
	if config.Logger.Development {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.OutputPaths = []string{
		"../e-commerce.log",
	}

	logger, err := cfg.Build()
	if err != nil {
		logger = zap.NewNop()
	}
	zap.ReplaceGlobals(logger)
}

func Close() {
	defer zap.L().Sync()
}
