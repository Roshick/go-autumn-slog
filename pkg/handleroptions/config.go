package handleroptions

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Roshick/go-autumn-slog/pkg/level"

	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

const (
	DefaultKeyLogLevel                = "LOG_LEVEL"
	DefaultKeyLogAttributeKeyMappings = "LOG_ATTRIBUTE_KEY_MAPPINGS"
)

type DefaultConfigImpl struct {
	VLogLevel                slog.Level
	VLogAttributeKeyMappings map[string]string
	VTimestampTransformer    func(time.Time) time.Time
}

func (c *DefaultConfigImpl) HandlerOptions() *slog.HandlerOptions {
	replaceAttr := func(_ []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey {
			attr.Value = slog.TimeValue(c.VTimestampTransformer(attr.Value.Time()))
		}
		if attr.Key == slog.LevelKey {
			logLevel := attr.Value.Any().(slog.Level)
			attr.Value = slog.StringValue(level.LevelToString(logLevel))

		}
		if mappedKey, ok := c.VLogAttributeKeyMappings[attr.Key]; ok {
			attr.Key = mappedKey
		}
		return attr
	}

	return &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
		Level:       c.VLogLevel,
	}
}

func DefaultConfigItems() []auconfigapi.ConfigItem {
	return []auconfigapi.ConfigItem{
		{
			Key:         DefaultKeyLogLevel,
			EnvName:     DefaultKeyLogLevel,
			Default:     "INFO",
			Description: "Minimum level of all logs.",
			Validate:    validateIsLogLevel,
		}, {
			Key:     DefaultKeyLogAttributeKeyMappings,
			EnvName: DefaultKeyLogAttributeKeyMappings,
			Default: "{}",
			Description: "Mappings for attribute keys of all logs. " +
				"Example: The entry [error: error.message] maps every attribute with key \"error\" to use the key \"error.message\" instead.",
			Validate: validateIsLogAttributeKeyMappings,
		},
	}
}

type ValuesProvider interface {
	ObtainValues(configItems []auconfigapi.ConfigItem) (map[string]string, error)
}

func ObtainDefaultConfig(provider ValuesProvider) (*DefaultConfigImpl, error) {
	values, err := provider.ObtainValues(DefaultConfigItems())
	if err != nil {
		return nil, fmt.Errorf("failed to obtain configuration values: %s", err.Error())
	}

	vLogLevel, _ := level.ParseLogLevel(values[DefaultKeyLogLevel])
	vLogAttributeKeyMappings, _ := parseLogAttributeKeyMappings(values[DefaultKeyLogAttributeKeyMappings])
	vTimestampTransformer := func(timestamp time.Time) time.Time {
		return timestamp.UTC()
	}
	return &DefaultConfigImpl{
		VLogLevel:                vLogLevel,
		VLogAttributeKeyMappings: vLogAttributeKeyMappings,
		VTimestampTransformer:    vTimestampTransformer,
	}, nil
}

func validateIsLogLevel(value string) error {
	_, err := level.ParseLogLevel(value)
	return err
}

func validateIsLogAttributeKeyMappings(value string) error {
	_, err := parseLogAttributeKeyMappings(value)
	return err
}

func parseLogAttributeKeyMappings(value string) (map[string]string, error) {
	var attributeKeyMappings map[string]string
	if err := json.Unmarshal([]byte(value), &attributeKeyMappings); err != nil {
		return nil, err
	}
	return attributeKeyMappings, nil
}
