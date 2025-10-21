package logger

import (
	"io"
	"log/slog"
)

var logLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func MustLoad(w io.Writer, level string) *slog.Logger {
	lvl, ok := logLevels[level]
	if !ok {
		panic("invalid logging level provided")
	}
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: lvl}))
}
