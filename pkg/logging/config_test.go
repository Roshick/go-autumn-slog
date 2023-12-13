package logging_test

import (
	"bytes"
	"github.com/Roshick/go-autumn-slog/pkg/level"
	"github.com/Roshick/go-autumn-slog/pkg/logging"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func getter(key string) string {
	values := map[string]string{
		logging.DefaultConfigKeyLevel:           "FATAL",
		logging.DefaultConfigKeyTimeTransformer: "ZERO",
		logging.DefaultConfigKeyAttributeKeyMappings: `
			{
				"time": "@timestamp"
			}
		`,
	}
	return values[key]
}

func TestObtainDefaultConfig_TextHandler(t *testing.T) {
	config := logging.NewConfig()

	err := config.ObtainValues(getter)
	assert.NoError(t, err)
	assert.Equal(t, level.Fatal, config.HandlerOptions().Level.Level())

	result := bytes.NewBuffer(nil)
	handler := slog.NewTextHandler(result, config.HandlerOptions())
	logger := slog.New(handler)
	logger.Log(nil, level.Panic, "this is a test")
	assert.Equal(t, "@timestamp=0001-01-01T00:00:00.000Z level=PANIC msg=\"this is a test\"\n", result.String())
}

func TestObtainDefaultConfig_JSONHandler(t *testing.T) {
	config := logging.NewConfig()

	err := config.ObtainValues(getter)
	assert.NoError(t, err)
	assert.Equal(t, level.Fatal, config.HandlerOptions().Level.Level())

	result := bytes.NewBuffer(nil)
	handler := slog.NewJSONHandler(result, config.HandlerOptions())
	logger := slog.New(handler)
	logger.Log(nil, level.Panic, "this is a test")
	assert.Equal(t, "{\"@timestamp\":\"0001-01-01T00:00:00Z\",\"level\":\"PANIC\",\"msg\":\"this is a test\"}\n", result.String())
}
