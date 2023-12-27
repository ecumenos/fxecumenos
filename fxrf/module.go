package fxrf

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewFactory),
)
