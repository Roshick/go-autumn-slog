package callbackhandler_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	callbackhandler "github.com/Roshick/go-autumn-slog/pkg/handlers/callback"
	"github.com/stretchr/testify/assert"
)

func newTestHandler(writer io.Writer, level slog.Level) slog.Handler {
	return slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				attr.Value = slog.TimeValue(time.Time{})
			}
			return attr
		},
	})
}

func TestHandler_RegisterContextCallback_AddsContextValueToAttribute(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelInfo)
	callbackHandler := callbackhandler.New(testHandler)
	err := callbackHandler.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
		record.Add("some-attr-key", ctx.Value("some-context-key"))
		return nil
	}, "test")
	assert.NoError(t, err)

	logger := slog.New(callbackHandler)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "some-context-key", "some-attr-value")
	logger.InfoContext(ctx, "this is a test")
	assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=INFO msg=\"this is a test\" some-attr-key=some-attr-value\n", result.String())
}

func TestHandler_RegisterContextCallback_ProducesErrorOnDuplicateCallbackRegistration(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelInfo)
	callbackHandler := callbackhandler.New(testHandler)
	err := callbackHandler.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
		return nil
	}, "test")
	assert.NoError(t, err)
	err = callbackHandler.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
		return nil
	}, "test")
	assert.Errorf(t, err, "failed to register callback with key test: callback with such key already exists")
}

func TestHandler_DeregisterContextCallback(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelInfo)
	callbackHandler := callbackhandler.New(testHandler)
	err := callbackHandler.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
		record.Add("some-attr-key", ctx.Value("some-context-key"))
		return nil
	}, "test")
	assert.NoError(t, err)
	callbackHandler.DeregisterContextCallback("test")

	logger := slog.New(callbackHandler)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "some-context-key", "some-attr-value")
	logger.InfoContext(ctx, "this is a test")
	assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=INFO msg=\"this is a test\"\n", result.String())
}

func TestHandler_Enabled(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelWarn)
	callbackHandler := callbackhandler.New(testHandler)

	logger := slog.New(callbackHandler)
	assert.False(t, logger.Enabled(context.TODO(), slog.LevelInfo))
}

func TestHandler_WithAttrs(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelInfo)
	callbackHandler := callbackhandler.New(testHandler)

	logger := slog.New(callbackHandler).With("some-attr-key", "some-attr-value")
	logger.Info("this is a test")
	assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=INFO msg=\"this is a test\" some-attr-key=some-attr-value\n", result.String())
}

func TestHandler_WithAttrs_InheritsIndependentCallbacks1(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelInfo)
	callbackHandler := callbackhandler.New(testHandler)
	err := callbackHandler.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
		record.Add("some-attr-key", "some-attr-value")
		return nil
	}, "test")
	assert.NoError(t, err)

	loggerParent := slog.New(callbackHandler)
	loggerChild := loggerParent.With("some-other-attr-key", "some-other-attr-value")
	callbackHandler.DeregisterContextCallback("test")
	loggerChild.Info("this is a test")
	assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=INFO msg=\"this is a test\" some-other-attr-key=some-other-attr-value some-attr-key=some-attr-value\n", result.String())
}

func TestHandler_WithAttrs_InheritsIndependentCallbacks2(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelInfo)
	callbackHandlerParent := callbackhandler.New(testHandler)

	loggerParent := slog.New(callbackHandlerParent)
	loggerChild := loggerParent.With("some-other-attr-key", "some-other-attr-value")

	callbackHandlerChild, ok := loggerChild.Handler().(*callbackhandler.Handler)
	assert.True(t, ok)
	err := callbackHandlerChild.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
		record.Add("some-attr-key", "some-attr-value")
		return nil
	}, "test")
	assert.NoError(t, err)

	loggerParent.Info("this is a test")
	assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=INFO msg=\"this is a test\"\n", result.String())
}

func TestHandler_WithGroup(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelInfo)
	callbackHandler := callbackhandler.New(testHandler)

	logger := slog.New(callbackHandler).WithGroup("some-group").With("some-attr-key", "some-attr-value")
	logger.Info("this is a test")
	assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=INFO msg=\"this is a test\" some-group.some-attr-key=some-attr-value\n", result.String())
}

func TestHandler_WithGroup_InheritsIndependentCallbacks1(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelInfo)
	callbackHandler := callbackhandler.New(testHandler)
	err := callbackHandler.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
		record.Add("some-attr-key", "some-attr-value")
		return nil
	}, "test")
	assert.NoError(t, err)

	loggerParent := slog.New(callbackHandler)
	loggerChild := loggerParent.WithGroup("some-group")
	callbackHandler.DeregisterContextCallback("test")
	loggerChild.Info("this is a test")
	assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=INFO msg=\"this is a test\" some-group.some-attr-key=some-attr-value\n", result.String())
}

func TestHandler_WithGroup_InheritsIndependentCallbacks2(t *testing.T) {
	result := bytes.NewBuffer(nil)
	testHandler := newTestHandler(result, slog.LevelInfo)
	callbackHandlerParent := callbackhandler.New(testHandler)

	loggerParent := slog.New(callbackHandlerParent)
	loggerChild := loggerParent.WithGroup("some-group")

	callbackHandlerChild, ok := loggerChild.Handler().(*callbackhandler.Handler)
	assert.True(t, ok)
	err := callbackHandlerChild.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
		record.Add("some-attr-key", "some-attr-value")
		return nil
	}, "test")
	assert.NoError(t, err)

	loggerParent.Info("this is a test")
	assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=INFO msg=\"this is a test\"\n", result.String())
}
