package telemetry

import (
	"io"
	"log/slog"
)

func NewLogger(w io.Writer, level, format string) *slog.Logger {
	var lvl slog.Level
	_ = lvl.UnmarshalText([]byte(level)) // debug, info, warn, error
	opts := &slog.HandlerOptions{Level: lvl}

	if format == "json" {
		return slog.New(slog.NewJSONHandler(w, opts))
	}

	return slog.New(slog.NewTextHandler(w, opts))
}
