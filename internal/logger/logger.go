package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"medods-test-task/internal/logger/config"
)

func NewLogger(cfg *config.Config) (*slog.Logger, error) {
	var writer io.Writer
	switch cfg.Output.Type {
	case "file":
		writer = &lumberjack.Logger{
			Filename:   cfg.Output.Path,
			MaxSize:    cfg.Output.MaxSize,
			MaxBackups: cfg.Output.MaxBackups,
			MaxAge:     cfg.Output.MaxAge,
			Compress:   cfg.Output.Compress,
		}
	case "stdout":
		fallthrough
	default:
		writer = os.Stdout
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: parseLogLevel(cfg.Level),
	}
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	return slog.New(handler), nil
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
