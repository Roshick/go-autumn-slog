package handleroptions

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Roshick/go-autumn-slog/pkg/level"

	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

const (
	DefaultKeyLogLevel                = "LOG_LEVEL"
	DefaultKeyLogAttributeKeyMappings = "LOG_ATTRIBUTE_KEY_MAPPINGS"
)

type TimestampTransformer func(time.Time) time.Time

type Config struct {
	vLogLevel                slog.Level
	vLogAttributeKeyMappings map[string]string
	vTimestampTransformer    TimestampTransformer
}

func NewDefaultConfig() *Config {
	return &Config{
		vTimestampTransformer: func(timestamp time.Time) time.Time {
			return timestamp.UTC()
		},
	}
}

func (c *Config) SetTimestampTransformer(transformer TimestampTransformer) {
	c.vTimestampTransformer = transformer
}

func (c *Config) HandlerOptions() *slog.HandlerOptions {
	replaceAttr := func(_ []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey {
			attr.Value = slog.TimeValue(c.vTimestampTransformer(attr.Value.Time()))
		}
		if attr.Key == slog.LevelKey {
			logLevel := attr.Value.Any().(slog.Level)
			attr.Value = slog.StringValue(level.LevelToString(logLevel))

		}
		if mappedKey, ok := c.vLogAttributeKeyMappings[attr.Key]; ok {
			attr.Key = mappedKey
		}
		return attr
	}

	return &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
		Level:       c.vLogLevel,
	}
}

func (c *Config) ConfigItems() []auconfigapi.ConfigItem {
	return []auconfigapi.ConfigItem{
		{
			Key:         DefaultKeyLogLevel,
			EnvName:     DefaultKeyLogLevel,
			Default:     "INFO",
			Description: "Minimum level of all logs.",
			Validate:    auconfigapi.ConfigNeedsNoValidation,
		}, {
			Key:     DefaultKeyLogAttributeKeyMappings,
			EnvName: DefaultKeyLogAttributeKeyMappings,
			Default: "{}",
			Description: "Mappings for attribute keys of all logs. " +
				"Example: The entry [error: error.message] maps every attribute with key \"error\" to use the key \"error.message\" instead.",
			Validate: auconfigapi.ConfigNeedsNoValidation,
		},
	}
}

func (c *Config) ObtainValues(getter func(string) string) error {
	if vLogLevel, err := level.ParseLogLevel(getter(DefaultKeyLogLevel)); err != nil {
		return err
	} else {
		c.vLogLevel = vLogLevel
	}

	if vLogAttributeKeyMappings, err := parseLogAttributeKeyMappings(getter(DefaultKeyLogAttributeKeyMappings)); err != nil {
		return err
	} else {
		c.vLogAttributeKeyMappings = vLogAttributeKeyMappings
	}

	return nil
}

func parseLogAttributeKeyMappings(value string) (map[string]string, error) {
	var attributeKeyMappings map[string]string
	if err := json.Unmarshal([]byte(value), &attributeKeyMappings); err != nil {
		return nil, err
	}
	return attributeKeyMappings, nil
}
