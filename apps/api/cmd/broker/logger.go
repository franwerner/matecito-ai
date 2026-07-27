package main

import (
	"io"
	"log/slog"

	"github.com/franwerner/matecito-ai/apps/api/internal/config"
)

// newLogger builds the broker's structured logger per observability/logging:
// text by default (readable when tailing the daemon), switchable to JSON,
// with a configurable minimum level.
func newLogger(cfg config.Config, w io.Writer) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(cfg.LogLevel)}

	var handler slog.Handler
	if cfg.LogFormat == "json" {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}
	return slog.New(handler)
}

// parseLevel defaults to info; config.Load already rejects any level string
// other than debug/info/warn/error before this runs, so the default branch
// here is just a defensive fallback, never expected to trigger.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
