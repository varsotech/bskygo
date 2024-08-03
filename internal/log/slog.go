package log

import (
	"log/slog"
)

// NewSlog creates a new default logger, using slog.
func NewSlog() Logger {
	return slog.Default()
}
