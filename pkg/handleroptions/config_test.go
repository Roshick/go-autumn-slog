package handleroptions_test

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/Roshick/go-autumn-slog/pkg/handleroptions"
	"github.com/Roshick/go-autumn-slog/pkg/level"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	"github.com/stretchr/testify/assert"
)

type testProvider struct{}

func (p *testProvider) ObtainValues(_ []auconfigapi.ConfigItem) (map[string]string, error) {
	return map[string]string{
		handleroptions.DefaultKeyLogLevel: "FATAL",
		handleroptions.DefaultKeyLogAttributeKeyMappings: `
			{
				"time": "@timestamp"
			}
		`,
	}, nil
}

func newTestHandler() handleroptions.ValuesProvider {
	return &testProvider{}
}

func TestObtainDefaultConfig_TextHandler(t *testing.T) {
	provider := newTestHandler()
	configuration, err := handleroptions.ObtainDefaultConfig(provider)
	assert.NoError(t, err)
	assert.Equal(t, level.Fatal, configuration.VLogLevel)

	configuration.VTimestampTransformer = func(timestamp time.Time) time.Time {
		return time.Time{}
	}
	result := bytes.NewBuffer(nil)
	handler := slog.NewTextHandler(result, configuration.HandlerOptions())
	logger := slog.New(handler)
	logger.Log(nil, level.Panic, "this is a test")
	assert.Equal(t, "@timestamp=0001-01-01T00:00:00.000Z level=PANIC msg=\"this is a test\"\n", result.String())
}

func TestObtainDefaultConfig_JSONHandler(t *testing.T) {
	provider := newTestHandler()
	configuration, err := handleroptions.ObtainDefaultConfig(provider)
	assert.NoError(t, err)
	assert.Equal(t, level.Fatal, configuration.VLogLevel)

	configuration.VTimestampTransformer = func(timestamp time.Time) time.Time {
		return time.Time{}
	}
	result := bytes.NewBuffer(nil)
	handler := slog.NewJSONHandler(result, configuration.HandlerOptions())
	logger := slog.New(handler)
	logger.Log(nil, level.Panic, "this is a test")
	assert.Equal(t, "{\"@timestamp\":\"0001-01-01T00:00:00Z\",\"level\":\"PANIC\",\"msg\":\"this is a test\"}\n", result.String())
}
