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

type LogLevel slog.Level

func ParseLogLevel(value string) (LogLevel, error) {
	mapping := map[string]LogLevel{
		"TRACE":  LogLevel(Trace),
		"DEBUG":  LogLevel(slog.LevelDebug),
		"INFO":   LogLevel(slog.LevelInfo),
		"WARN":   LogLevel(slog.LevelWarn),
		"ERROR":  LogLevel(slog.LevelError),
		"FATAL":  LogLevel(Fatal),
		"PANIC":  LogLevel(Panic),
		"SILENT": LogLevel(Silent),
	}
	if logLevel, ok := mapping[value]; ok {
		return logLevel, nil
	}
	mappingKeys := make([]string, 0)
	for key := range mapping {
		mappingKeys = append(mappingKeys, key)
	}
	return LogLevel(Silent), fmt.Errorf("failed to parse log level %s, must be one of %v", value, mappingKeys)
}

func LogLevelToString(level LogLevel) string {
	mapping := map[LogLevel]string{
		LogLevel(Trace):           "TRACE",
		LogLevel(slog.LevelDebug): "DEBUG",
		LogLevel(slog.LevelInfo):  "INFO",
		LogLevel(slog.LevelWarn):  "WARN",
		LogLevel(slog.LevelError): "ERROR",
		LogLevel(Fatal):           "FATAL",
		LogLevel(Panic):           "PANIC",
		LogLevel(Silent):          "SILENT",
	}
	if value, ok := mapping[level]; ok {
		return value
	}
	return "UNKNOWN"
}
