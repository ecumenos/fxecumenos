package loggerfx

import (
	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/fxecumenos/loggerfx/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type loggerParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	ServiceName fxecumenos.ServiceName
	Config      *Config
}

type Config struct {
	Prod bool
}

var Module = fx.Options(
	fx.Provide(
		func(params loggerParams) (*zap.Logger, error) {
			return logger.NewZapLogger(params.ServiceName, params.Config.Prod, params.Lifecycle)
		},
		logger.ZapSugared,
	),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger}
	}),
)
