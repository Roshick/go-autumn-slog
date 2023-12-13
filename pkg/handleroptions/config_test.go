package handleroptions_test

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/Roshick/go-autumn-slog/pkg/handleroptions"
	"github.com/Roshick/go-autumn-slog/pkg/level"
	"github.com/stretchr/testify/assert"
)

func getter(key string) string {
	values := map[string]string{
		handleroptions.DefaultKeyLogLevel: "FATAL",
		handleroptions.DefaultKeyLogAttributeKeyMappings: `
			{
				"time": "@timestamp"
			}
		`,
	}
	return values[key]
}

func TestObtainDefaultConfig_TextHandler(t *testing.T) {
	config := handleroptions.NewDefaultConfig()

	err := config.ObtainValues(getter)
	assert.NoError(t, err)
	assert.Equal(t, level.Fatal, config.HandlerOptions().Level.Level())

	config.SetTimestampTransformer(func(timestamp time.Time) time.Time {
		return time.Time{}
	})
	result := bytes.NewBuffer(nil)
	handler := slog.NewTextHandler(result, config.HandlerOptions())
	logger := slog.New(handler)
	logger.Log(nil, level.Panic, "this is a test")
	assert.Equal(t, "@timestamp=0001-01-01T00:00:00.000Z level=PANIC msg=\"this is a test\"\n", result.String())
}

func TestObtainDefaultConfig_JSONHandler(t *testing.T) {
	config := handleroptions.NewDefaultConfig()

	err := config.ObtainValues(getter)
	assert.NoError(t, err)
	assert.Equal(t, level.Fatal, config.HandlerOptions().Level.Level())

	config.SetTimestampTransformer(func(timestamp time.Time) time.Time {
		return time.Time{}
	})
	result := bytes.NewBuffer(nil)
	handler := slog.NewJSONHandler(result, config.HandlerOptions())
	logger := slog.New(handler)
	logger.Log(nil, level.Panic, "this is a test")
	assert.Equal(t, "{\"@timestamp\":\"0001-01-01T00:00:00Z\",\"level\":\"PANIC\",\"msg\":\"this is a test\"}\n", result.String())
}
