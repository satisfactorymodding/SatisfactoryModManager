package app

import (
	"log/slog"

	"github.com/spf13/viper"
)

func (a *app) GetDebug() bool {
	return viper.GetString("log") == "debug"
}

func (a *app) SetDebug(debug bool) {
	slog.Info("Setting debug mode", slog.Bool("debug", debug))
	if debug {
		viper.Set("log", "debug")
	} else {
		viper.Set("log", "info")
	}
}
