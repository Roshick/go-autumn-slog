package logging

import (
	"encoding/json"
	"fmt"
	"github.com/Roshick/go-autumn-slog/pkg/level"
	"log/slog"
	"time"

	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

const (
	DefaultConfigKeyLevel                = "LOG_LEVEL"
	DefaultConfigKeyTimeTransformer      = "LOG_TIME_TRANSFORMER"
	DefaultConfigKeyAttributeKeyMappings = "LOG_ATTRIBUTE_KEY_MAPPINGS"
)

type TimeTransformer func(time.Time) time.Time

type Config struct {
	vLogLevel                slog.Level
	vLogAttributeKeyMappings map[string]string
	vTimeTransformer         TimeTransformer
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LogLevel() slog.Level {
	return c.vLogLevel
}

func (c *Config) HandlerOptions() *slog.HandlerOptions {
	replaceAttr := func(_ []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey {
			attr.Value = slog.TimeValue(c.vTimeTransformer(attr.Value.Time()))
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

const DefaultMapping = `{
  "time": "@timestamp",
  "level": "log.level",
  "msg": "message",
  "error": "error.message"
}`

func (c *Config) ConfigItems() []auconfigapi.ConfigItem {
	return []auconfigapi.ConfigItem{
		{
			Key:     DefaultConfigKeyLevel,
			EnvName: DefaultConfigKeyLevel,
			Default: "INFO",
			Description: "Minimum level of all logs. \n" +
				"Supported values: TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC, SILENT",
			Validate: auconfigapi.ConfigNeedsNoValidation,
		}, {
			Key:     DefaultConfigKeyTimeTransformer,
			EnvName: DefaultConfigKeyTimeTransformer,
			Default: "UTC",
			Description: "Type of transformation applied to each record's timestamp. Useful for testing purposes. \n" +
				"Supported values: UTC, ZERO",
			Validate: auconfigapi.ConfigNeedsNoValidation,
		}, {
			Key:     DefaultConfigKeyAttributeKeyMappings,
			EnvName: DefaultConfigKeyAttributeKeyMappings,
			Default: DefaultMapping,
			Description: "Mappings for attribute keys of all logs. \n" +
				"Example: The entry [error: error.message] maps every attribute with key \"error\" to use the key \"error.message\" instead.",
			Validate: auconfigapi.ConfigNeedsNoValidation,
		},
	}
}

func (c *Config) ObtainValues(getter func(string) string) error {
	if vLogLevel, err := level.ParseLogLevel(getter(DefaultConfigKeyLevel)); err != nil {
		return err
	} else {
		c.vLogLevel = vLogLevel
	}

	if vTimeTransformer, err := parseTimeTransformer(getter(DefaultConfigKeyTimeTransformer)); err != nil {
		return err
	} else {
		c.vTimeTransformer = vTimeTransformer
	}

	if vLogAttributeKeyMappings, err := parseLogAttributeKeyMappings(getter(DefaultConfigKeyAttributeKeyMappings)); err != nil {
		return err
	} else {
		c.vLogAttributeKeyMappings = vLogAttributeKeyMappings
	}

	return nil
}

func parseTimeTransformer(value string) (TimeTransformer, error) {
	switch value {
	case "UTC":
		return func(timestamp time.Time) time.Time {
			return timestamp.UTC()
		}, nil
	case "ZERO":
		return func(_ time.Time) time.Time {
			return time.Time{}
		}, nil
	default:
		return nil, fmt.Errorf("invalid time transformer: '%s'", value)
	}
}

func parseLogAttributeKeyMappings(value string) (map[string]string, error) {
	var attributeKeyMappings map[string]string
	if err := json.Unmarshal([]byte(value), &attributeKeyMappings); err != nil {
		return nil, err
	}
	return attributeKeyMappings, nil
}
