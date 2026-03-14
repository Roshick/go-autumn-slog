package logging_test

import (
	"bytes"
	"os"
	"testing"

	"log/slog"

	"github.com/Roshick/go-autumn-slog"
	"github.com/Roshick/go-autumn-slog/level"
	"github.com/stretchr/testify/assert"
)

func TestObtainDefaultConfig_TextHandler(t *testing.T) {
	// Set env vars
	os.Setenv("LOG_LEVEL", "FATAL")
	os.Setenv("LOG_TIME_TRANSFORMER", "ZERO")
	os.Setenv("LOG_ATTRIBUTE_KEY_MAPPINGS", `{"time": "@timestamp"}`)
	defer func() {
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_TIME_TRANSFORMER")
		os.Unsetenv("LOG_ATTRIBUTE_KEY_MAPPINGS")
	}()

	config := logging.NewConfig()

	err := config.ObtainValuesFromEnv()
	assert.NoError(t, err)
	assert.Equal(t, level.Fatal, config.HandlerOptions().Level)

	result := bytes.NewBuffer(nil)
	handler := slog.NewTextHandler(result, config.HandlerOptions())
	logger := slog.New(handler)
	logger.Log(nil, level.Panic, "this is a test")
	assert.Equal(t, "@timestamp=0001-01-01T00:00:00.000Z level=PANIC msg=\"this is a test\"\n", result.String())
}

func TestObtainDefaultConfig_JSONHandler(t *testing.T) {
	// Set env vars
	os.Setenv("LOG_LEVEL", "FATAL")
	os.Setenv("LOG_TIME_TRANSFORMER", "ZERO")
	os.Setenv("LOG_ATTRIBUTE_KEY_MAPPINGS", `{"time": "@timestamp"}`)
	defer func() {
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_TIME_TRANSFORMER")
		os.Unsetenv("LOG_ATTRIBUTE_KEY_MAPPINGS")
	}()

	config := logging.NewConfig()

	err := config.ObtainValuesFromEnv()
	assert.NoError(t, err)
	assert.Equal(t, level.Fatal, config.HandlerOptions().Level)

	result := bytes.NewBuffer(nil)
	handler := slog.NewJSONHandler(result, config.HandlerOptions())
	logger := slog.New(handler)
	logger.Log(nil, level.Panic, "this is a test")
	assert.Equal(t, "{\"@timestamp\":\"0001-01-01T00:00:00Z\",\"level\":\"PANIC\",\"msg\":\"this is a test\"}\n", result.String())
}
