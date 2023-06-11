package main

import "github.com/rs/zerolog/log"

type wailsZeroLogLogger struct{}

func (l wailsZeroLogLogger) Print(message string) {
	log.Trace().Msg(message)
}

func (l wailsZeroLogLogger) Trace(message string) {
	log.Trace().Msg(message)
}

func (l wailsZeroLogLogger) Debug(message string) {
	log.Debug().Msg(message)
}

func (l wailsZeroLogLogger) Info(message string) {
	log.Info().Msg(message)
}

func (l wailsZeroLogLogger) Warning(message string) {
	log.Warn().Msg(message)
}

func (l wailsZeroLogLogger) Error(message string) {
	log.Error().Msg(message)
}

func (l wailsZeroLogLogger) Fatal(message string) {
	log.Fatal().Msg(message)
}
