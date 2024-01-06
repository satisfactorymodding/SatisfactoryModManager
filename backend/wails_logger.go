package backend

import (
	"log/slog"
	"os"
)

type WailsZeroLogLogger struct{}

func (l WailsZeroLogLogger) Print(message string) {
	slog.Debug(message)
}

func (l WailsZeroLogLogger) Trace(message string) {
	slog.Debug(message)
}

func (l WailsZeroLogLogger) Debug(message string) {
	slog.Debug(message)
}

func (l WailsZeroLogLogger) Info(message string) {
	slog.Info(message)
}

func (l WailsZeroLogLogger) Warning(message string) {
	slog.Warn(message)
}

func (l WailsZeroLogLogger) Error(message string) {
	slog.Error(message)
}

func (l WailsZeroLogLogger) Fatal(message string) {
	slog.Error(message)
	os.Exit(1)
}
