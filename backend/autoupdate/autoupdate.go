package autoupdate

import (
	"log/slog"
	"time"

	"github.com/spf13/viper"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/source/github"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater"
)

var Updater *updater.Updater

var checkStarted bool

func Init() {
	config := updater.Config{
		Source:         github.MakeGithubProvider(viper.GetString("github-release-repo")),
		CurrentVersion: viper.GetString("version"),
	}
	updateType := getUpdateType()
	// Some builds cannot (or should not) auto-update
	if updateType != nil {
		config.File = updateType.ArtifactName
		config.Apply = updateType.Apply
	}
	Updater = updater.MakeUpdater(config)
}

var (
	updateCheckTicker *time.Ticker
	updateCheckStop   = make(chan bool)
)

func CheckInterval(interval time.Duration) {
	if checkStarted {
		return
	}
	checkStarted = true
	updateCheckTicker = time.NewTicker(interval)
	updateCheckStop = make(chan bool)
	go func() {
		err := Updater.CheckForUpdate()
		if err != nil {
			slog.Error("failed to check for update", slog.Any("error", err))
		}
		for range updateCheckTicker.C {
			select {
			case <-updateCheckStop:
				return
			case <-updateCheckTicker.C:
				err := Updater.CheckForUpdate()
				if err != nil {
					slog.Error("failed to check for update", slog.Any("error", err))
				}
			}
		}
	}()
}

func OnExit(restart bool) error {
	close(updateCheckStop)
	if Updater == nil {
		// No updater for this build type
		return nil
	}
	return Updater.OnExit(restart)
}
