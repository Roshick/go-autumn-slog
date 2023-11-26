package level

import (
	"fmt"
	"log/slog"
)

const (
	Trace  = slog.Level(-8)
	Fatal  = slog.Level(12)
	Panic  = slog.Level(16)
	Silent = slog.Level(20)
)

func ParseLogLevel(value string) (slog.Level, error) {
	mapping := map[string]slog.Level{
		"TRACE":  Trace,
		"DEBUG":  slog.LevelDebug,
		"INFO":   slog.LevelInfo,
		"WARN":   slog.LevelWarn,
		"ERROR":  slog.LevelError,
		"FATAL":  Fatal,
		"PANIC":  Panic,
		"SILENT": Silent,
	}
	if logLevel, ok := mapping[value]; ok {
		return logLevel, nil
	}
	mappingKeys := make([]string, 0)
	for key := range mapping {
		mappingKeys = append(mappingKeys, key)
	}
	return Silent, fmt.Errorf("failed to parse log level %s, must be one of %v", value, mappingKeys)
}

func LevelToString(level slog.Level) string {
	mapping := map[slog.Level]string{
		Trace:           "TRACE",
		slog.LevelDebug: "DEBUG",
		slog.LevelInfo:  "INFO",
		slog.LevelWarn:  "WARN",
		slog.LevelError: "ERROR",
		Fatal:           "FATAL",
		Panic:           "PANIC",
		Silent:          "SILENT",
	}
	if value, ok := mapping[level]; ok {
		return value
	}
	return "UNKNOWN"
}
