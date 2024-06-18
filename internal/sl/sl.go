package sl

import (
	"log/slog"
)

// Err синтаксический сахар для slog
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
