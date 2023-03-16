package wailslogging

import "github.com/rs/zerolog/log"

type WailsZeroLogLogger struct{}

func (l WailsZeroLogLogger) Print(message string) {
	log.Trace().Msg(message)
}

func (l WailsZeroLogLogger) Trace(message string) {
	log.Trace().Msg(message)
}

func (l WailsZeroLogLogger) Debug(message string) {
	log.Debug().Msg(message)
}

func (l WailsZeroLogLogger) Info(message string) {
	log.Info().Msg(message)
}

func (l WailsZeroLogLogger) Warning(message string) {
	log.Warn().Msg(message)
}

func (l WailsZeroLogLogger) Error(message string) {
	log.Error().Msg(message)
}

func (l WailsZeroLogLogger) Fatal(message string) {
	log.Fatal().Msg(message)
}
