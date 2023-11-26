package callbackhandler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type Handler struct {
	wrappedHandler slog.Handler
	callbacks      *sync.Map
}

type CallbackFunc func(ctx context.Context, record *slog.Record) error

func New(wrappedHandler slog.Handler) *Handler {
	var callbacks sync.Map
	return &Handler{
		wrappedHandler: wrappedHandler,
		callbacks:      &callbacks,
	}
}

func (h *Handler) RegisterContextCallback(callback CallbackFunc, callbackKey string) error {
	if _, ok := h.callbacks.Load(callbackKey); ok {
		return fmt.Errorf("failed to register callback with key %s: callback with such key already exists", callbackKey)
	}
	h.callbacks.Store(callbackKey, callback)
	return nil
}

func (h *Handler) DeregisterContextCallback(callbackKey string) {
	h.callbacks.Delete(callbackKey)
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.wrappedHandler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var callbacksCopy sync.Map
	h.callbacks.Range(func(key any, value any) bool {
		callbacksCopy.Store(key, value)
		return true
	})
	return &Handler{
		wrappedHandler: h.wrappedHandler.WithAttrs(attrs),
		callbacks:      &callbacksCopy,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	var callbacksCopy sync.Map
	h.callbacks.Range(func(key any, value any) bool {
		callbacksCopy.Store(key, value)
		return true
	})
	return &Handler{
		wrappedHandler: h.wrappedHandler.WithGroup(name),
		callbacks:      &callbacksCopy,
	}
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	var err error
	h.callbacks.Range(func(_ any, value any) bool {
		callback, ok := value.(CallbackFunc)
		if !ok {
			err = fmt.Errorf("failed to retrieve context callback")
			return false
		}
		if err = callback(ctx, &record); err != nil {
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	return h.wrappedHandler.Handle(ctx, record)
}
