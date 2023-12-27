package fxrf

import (
	"net/http"

	"go.uber.org/zap"
)

//go:generate mockery --name=Factory

type Factory interface {
	NewWriter(rw http.ResponseWriter) Writer
}

type Config struct {
	WriteLogs bool
}

func NewFactory(l *zap.Logger, cfg *Config) Factory {
	return &factory{
		l:         l,
		writeLogs: cfg.WriteLogs,
	}
}

type factory struct {
	l         *zap.Logger
	writeLogs bool
}

func (f *factory) NewWriter(rw http.ResponseWriter) Writer {
	return NewWriter(f.l, rw, f.writeLogs)
}
