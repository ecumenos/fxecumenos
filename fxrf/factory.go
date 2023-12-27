package fxrf

import (
	"net/http"

	"github.com/ecumenos/fxecumenos"
	"go.uber.org/zap"
)

//go:generate mockery --name=Factory

type Factory interface {
	NewWriter(rw http.ResponseWriter) Writer
}

type Config struct {
	WriteLogs bool
}

func NewFactory(l *zap.Logger, cfg *Config, version fxecumenos.Version) Factory {
	return &factory{
		l:         l,
		writeLogs: cfg.WriteLogs,
		version:   version,
	}
}

type factory struct {
	l         *zap.Logger
	writeLogs bool
	version   fxecumenos.Version
}

func (f *factory) NewWriter(rw http.ResponseWriter) Writer {
	return NewWriter(f.l, rw, f.version, f.writeLogs)
}
