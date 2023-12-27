package logger

import (
	"context"

	"github.com/ecumenos/fxecumenos"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewZapLogger(serviceName fxecumenos.ServiceName, prod bool, lc fx.Lifecycle) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	if prod {
		logger, err = NewProductionLogger(string(serviceName))
	} else {
		logger, err = NewDevelopmentLogger(string(serviceName))
	}
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)

	lc.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			_ = logger.Sync()
			return nil
		},
	})

	return logger, nil
}

func ZapSugared(log *zap.Logger) *zap.SugaredLogger {
	return log.Sugar()
}
