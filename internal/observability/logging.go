package observability

import (
	"io"
	"log/slog"
	"os"
	"time"
)

// NewJSONLogger builds a stderr JSON logger.
func NewJSONLogger(level slog.Level) *slog.Logger {
	handler := NewJSONHandler(os.Stderr, level)
	return slog.New(handler)
}

// NewJSONHandler builds a JSON slog handler using second-level UTC timestamps
// formatted as time.DateTime.
func NewJSONHandler(w io.Writer, level slog.Level) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key != slog.TimeKey {
				return attr
			}

			t, ok := attr.Value.Any().(time.Time)
			if !ok {
				return attr
			}
			return slog.String(slog.TimeKey, t.UTC().Truncate(time.Second).Format(time.DateTime))
		},
	})
}
