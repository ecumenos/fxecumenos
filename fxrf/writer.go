package fxrf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/go-toolkit/contextutils"
	"github.com/ecumenos/go-toolkit/httputils"
	"github.com/ecumenos/go-toolkit/timeutils"
	"go.uber.org/zap"
)

//go:generate mockery --name=Writer

type Writer interface {
	SetLogger(logger *zap.Logger)
	WriteSuccess(ctx context.Context, payload interface{}, opts ...ResponseBuildOption) error
	WriteFail(ctx context.Context, data interface{}, opts ...ResponseBuildOption) error
	WriteError(ctx context.Context, msg string, cause error, opts ...ResponseBuildOption) error
}

type writer struct {
	rw         http.ResponseWriter
	l          *zap.Logger
	writeLogs  bool
	appVersion fxecumenos.Version
}

func NewWriter(l *zap.Logger, rw http.ResponseWriter, appVersion fxecumenos.Version, writeLogs bool) Writer {
	return &writer{
		rw:         rw,
		l:          l,
		writeLogs:  writeLogs,
		appVersion: appVersion,
	}
}

func (w *writer) write(payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if _, err := w.rw.Write(b); err != nil {
		return err
	}

	return nil
}

func (w *writer) writeHeaders(headers map[string]string, statusCode int) {
	w.rw.Header().Set("Content-Type", "application/json")
	for key, value := range headers {
		w.rw.Header().Set(key, value)
	}
	w.rw.WriteHeader(statusCode)
}

func (w *writer) SetLogger(logger *zap.Logger) {
	w.l = logger
}

type Status string

const (
	SuccessStatus Status = "success"
	FailureStatus Status = "failure"
	ErrorStatus   Status = "error"
)

type SuccessResp[T interface{}] struct {
	Status   Status    `json:"status"`
	Data     T         `json:"data"`
	Metadata *Metadata `json:"metadata"`
}

func (w *writer) WriteSuccess(ctx context.Context, payload interface{}, opts ...ResponseBuildOption) error {
	metadata, err := w.getMetadata(ctx)
	if err != nil {
		return err
	}
	rb := &responseBuilder{
		httpStatusCode: http.StatusOK,
		data:           payload,
		l:              w.l,
	}
	for _, opt := range opts {
		opt(rb)
	}
	if rb.httpStatusCode < http.StatusOK || rb.httpStatusCode > 299 {
		return fmt.Errorf("success response must have status code in range 200..299 (status code = %v)", rb.httpStatusCode)
	}
	if w.writeLogs {
		w.l.Info("responding success response", zap.Any("data", payload), zap.Int("status_code", rb.httpStatusCode))
	}

	w.writeHeaders(nil, rb.httpStatusCode)
	return w.write(&SuccessResp[interface{}]{
		Data:     rb.data,
		Metadata: metadata,
		Status:   SuccessStatus,
	})
}

type FailureResp[T interface{}] struct {
	Status   Status    `json:"status"`
	Data     T         `json:"data"`
	Message  string    `json:"message"`
	Metadata *Metadata `json:"metadata"`
}

func (w *writer) WriteFail(ctx context.Context, data interface{}, opts ...ResponseBuildOption) error {
	metadata, err := w.getMetadata(ctx)
	if err != nil {
		return err
	}
	rb := &responseBuilder{
		httpStatusCode: http.StatusBadRequest,
		data:           data,
		l:              w.l,
	}
	for _, opt := range opts {
		opt(rb)
	}
	if rb.httpStatusCode < http.StatusBadRequest || rb.httpStatusCode > 499 {
		return fmt.Errorf("fail response must have status code in range 400..499 (status code = %v)", rb.httpStatusCode)
	}
	if w.writeLogs {
		w.l.Info("responding fail response", zap.Any("data", rb.data), zap.Error(rb.cause),
			zap.Int("status_code", rb.httpStatusCode), zap.String("msg", rb.message))
	}

	w.writeHeaders(nil, http.StatusBadRequest)
	return w.write(&FailureResp[interface{}]{
		Data:     rb.data,
		Message:  rb.message,
		Metadata: metadata,
		Status:   FailureStatus,
	})
}

type ErrorResp struct {
	Status   Status    `json:"status"`
	Message  string    `json:"message"`
	Metadata *Metadata `json:"metadata"`
}

func (w *writer) WriteError(ctx context.Context, msg string, cause error, opts ...ResponseBuildOption) error {
	metadata, err := w.getMetadata(ctx)
	if err != nil {
		return err
	}
	rb := &responseBuilder{
		httpStatusCode: http.StatusInternalServerError,
		message:        msg,
		cause:          cause,
		l:              w.l,
	}
	for _, opt := range opts {
		opt(rb)
	}
	if rb.httpStatusCode < http.StatusInternalServerError || rb.httpStatusCode > 599 {
		return fmt.Errorf("error response must have status code in range 500..599 (status code = %v)", rb.httpStatusCode)
	}
	if w.writeLogs {
		w.l.Info("responding error response", zap.Error(rb.cause), zap.String("msg", rb.message), zap.Int("status_code", rb.httpStatusCode))
	}

	w.writeHeaders(nil, rb.httpStatusCode)
	return w.write(&ErrorResp{
		Message:  rb.message,
		Metadata: metadata,
		Status:   ErrorStatus,
	})
}

type responseBuilder struct {
	httpStatusCode int
	message        string
	cause          error
	data           interface{}
	l              *zap.Logger
}

type ResponseBuildOption func(b *responseBuilder)

func WithHTTPStatusCode(code int) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.httpStatusCode = code
	}
}

func WithMessage(msg string) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.message = msg
	}
}

func WithCause(err error) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.cause = err
	}
}

func WithData(data interface{}) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.data = data
	}
}

func WithLogger(l *zap.Logger) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.l = l
	}
}

type Metadata struct {
	RequestID string `json:"requestId"`
	Duration  int    `json:"duration"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

func (w *writer) getMetadata(ctx context.Context) (*Metadata, error) {
	duration, err := httputils.GetRequestDuration(ctx)
	if err != nil {
		return nil, err
	}

	return &Metadata{
		RequestID: contextutils.GetValueFromContext(ctx, contextutils.RequestIDKey),
		Timestamp: timeutils.TimeToString(time.Now()),
		Duration:  duration,
		Version:   strings.Join(w.appVersion.Build, "."),
	}, nil
}
