package logger

import "log/slog"

func NewLogger() *slog.Logger {
	log := slog.Default()

	return log
}
