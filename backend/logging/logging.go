package logging

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
)

func Init() {
	handlers := make([]slog.Handler, 0)

	if _, err := os.Stdout.Stat(); err == nil {
		// Only add the stdout handler if it is writable.
		// Otherwise, the fanout handler would have the first handler error,
		// and will not get to use the file handler.
		handlers = append(handlers, tint.NewHandler(os.Stdout, &tint.Options{
			Level:      settingsLogLevel{},
			AddSource:  true,
			TimeFormat: time.RFC3339,
		}))
	}

	if viper.GetString("log-file") != "" {
		logFile := &lumberjack.Logger{
			Filename:   viper.GetString("log-file"),
			MaxSize:    10, // megabytes
			MaxBackups: 5,
			MaxAge:     30, // days
		}

		handlers = append(handlers, slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: settingsLogLevel{},
		}))
	}

	slog.SetDefault(
		slog.New(
			slogmulti.
				Pipe(newRedactGamePathCredentialsMiddleware()).
				Handler(slogmulti.Fanout(handlers...)),
		),
	)
}

type settingsLogLevel struct{}

func (v settingsLogLevel) Level() slog.Level {
	if settings.Settings.Debug {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}
