package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	noophandler "github.com/Roshick/go-autumn-slog/handlers/noop"
	"github.com/Roshick/go-autumn-slog/level"

	auloggingapi "github.com/StephanHCB/go-autumn-logging/api"
)

const (
	ErrorKey = "error"
)

func NewLoggerFromEnv() (*slog.Logger, error) {
	config := NewConfig()
	if err := config.ObtainValuesFromEnv(); err != nil {
		return nil, fmt.Errorf("failed to obtain logger config values from environment: %w", err)
	}

	if config.LogStyle == LogStyleJSON {
		return slog.New(slog.NewJSONHandler(os.Stderr, config.HandlerOptions())), nil
	}
	return slog.New(slog.NewTextHandler(os.Stderr, config.HandlerOptions())), nil
}

type Logging struct {
	logger *slog.Logger
}

func New() *Logging {
	return &Logging{}
}

func (l *Logging) WithLogger(logger *slog.Logger) *Logging {
	return &Logging{
		logger: logger,
	}
}

func (l *Logging) Logger() *slog.Logger {
	currentLogger := l.logger
	if currentLogger == nil {
		currentLogger = slog.Default()
	}
	if currentLogger == nil {
		currentLogger = slog.New(noophandler.New())
	}
	return currentLogger
}

func (l *Logging) Ctx(ctx context.Context) auloggingapi.ContextAwareLoggingImplementation {
	currentLogger := FromContext(ctx)
	if currentLogger == nil {
		currentLogger = l.Logger()
	}
	return &withContext{
		Logging: l.WithLogger(currentLogger),
		ctx:     ctx,
	}
}

func (l *Logging) NoCtx() auloggingapi.ContextAwareLoggingImplementation {
	currentLogger := l.Logger()
	return &withContext{
		Logging: l.WithLogger(currentLogger),
		ctx:     nil,
	}
}

// withContext

type withContext struct {
	*Logging
	ctx context.Context
}

func (w *withContext) Trace() auloggingapi.LeveledLoggingImplementation {
	return &withLevel{
		withContext: w,
		level:       level.Trace,
	}
}

func (w *withContext) Debug() auloggingapi.LeveledLoggingImplementation {
	return &withLevel{
		withContext: w,
		level:       slog.LevelDebug,
	}
}

func (w *withContext) Info() auloggingapi.LeveledLoggingImplementation {
	return &withLevel{
		withContext: w,
		level:       slog.LevelInfo,
	}
}

func (w *withContext) Warn() auloggingapi.LeveledLoggingImplementation {
	return &withLevel{
		withContext: w,
		level:       slog.LevelWarn,
	}
}

func (w *withContext) Error() auloggingapi.LeveledLoggingImplementation {
	return &withLevel{
		withContext: w,
		level:       slog.LevelError,
	}
}

func (w *withContext) Fatal() auloggingapi.LeveledLoggingImplementation {
	return &withLevel{
		withContext: w,
		level:       level.Fatal,
	}
}

func (w *withContext) Panic() auloggingapi.LeveledLoggingImplementation {
	return &withLevel{
		withContext: w,
		level:       level.Panic,
	}
}

// withLevel

type withLevel struct {
	*withContext
	level slog.Level
}

func (w *withLevel) WithErr(err error) auloggingapi.LeveledLoggingImplementation {
	if err != nil {
		w.logger = w.logger.With(ErrorKey, err.Error())
	}
	return w
}

func (w *withLevel) With(key string, value string) auloggingapi.LeveledLoggingImplementation {
	w.logger = w.logger.With(key, value)
	return w
}

func (w *withLevel) Print(v ...any) {
	message := fmt.Sprint(v...)
	w.produceLog(message)
}

func (w *withLevel) Printf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	w.produceLog(message)
}

func (w *withLevel) produceLog(message string) {
	w.logger.Log(w.ctx, w.level, message)
}

type contextKey struct{}

func FromContext(ctx context.Context) *slog.Logger {
	if value := ctx.Value(contextKey{}); value != nil {
		return value.(*slog.Logger)
	}
	return nil
}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}
