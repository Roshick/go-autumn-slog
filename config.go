package logging

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/Roshick/go-autumn-slog/level"
	"github.com/caarlos0/env/v11"
)

type LogType int64

const (
	LogStylePlain LogType = iota
	LogStyleJSON
)

type TimeTransformer func(time.Time) time.Time

type Config struct {
	LogStyle                LogType           `env:"LOG_STYLE"                  envDefault:"PLAIN"`
	LogLevel                level.LogLevel    `env:"LOG_LEVEL"                  envDefault:"INFO"`
	LogAttributeKeyMappings map[string]string `env:"LOG_ATTRIBUTE_KEY_MAPPINGS" envDefault:"{\"time\":\"@timestamp\",\"level\":\"log.level\",\"msg\":\"message\",\"error\":\"error.message\"}"`
	TimeTransformer         TimeTransformer   `env:"LOG_TIME_TRANSFORMER" envDefault:"UTC"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) HandlerOptions() *slog.HandlerOptions {
	replaceAttr := func(_ []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey {
			attr.Value = slog.TimeValue(c.TimeTransformer(attr.Value.Time()))
		}
		if attr.Key == slog.LevelKey {
			logLevel := attr.Value.Any().(slog.Level)
			attr.Value = slog.StringValue(level.LogLevelToString(level.LogLevel(logLevel)))
		}
		if mappedKey, ok := c.LogAttributeKeyMappings[attr.Key]; ok {
			attr.Key = mappedKey
		}
		return attr
	}

	return &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
		Level:       slog.Level(c.LogLevel),
	}
}

func (c *Config) ObtainValuesFromEnv() error {
	return env.ParseWithOptions(c, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(LogType(0)): func(v string) (any, error) {
				return parseLogType(v)
			},
			reflect.TypeOf(level.LogLevel(0)): func(v string) (any, error) {
				return level.ParseLogLevel(v)
			},
			reflect.TypeOf(map[string]string{}): func(v string) (any, error) {
				return parseLogAttributeKeyMappings(v)
			},
			reflect.TypeOf(TimeTransformer(nil)): func(v string) (any, error) {
				return parseTimeTransformer(v)
			},
		},
	})
}

func parseLogType(value string) (LogType, error) {
	switch value {
	case "JSON":
		return LogStyleJSON, nil
	case "PLAIN":
		return LogStylePlain, nil
	default:
		return LogType(0), fmt.Errorf("invalid log type: '%s'", value)
	}
}

func parseLogAttributeKeyMappings(value string) (map[string]string, error) {
	var attributeKeyMappings map[string]string
	if err := json.Unmarshal([]byte(value), &attributeKeyMappings); err != nil {
		return nil, err
	}
	return attributeKeyMappings, nil
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
