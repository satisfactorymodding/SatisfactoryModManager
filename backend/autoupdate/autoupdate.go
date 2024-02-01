package autoupdate

import (
	"log/slog"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/source/github"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater"
)

var Updater *updater.Updater

var checkStarted bool

func Init() {
	if Updater != nil {
		return
	}
	Updater = updater.MakeUpdater(makeUpdaterConfig())
}

func makeUpdaterConfig() updater.Config {
	currentVersion, err := semver.NewVersion(viper.GetString("version"))
	if err != nil {
		if shouldUseUpdater() {
			slog.Error("failed to parse current version, using 0.0.0-unknown", slog.Any("error", err))
		}
		currentVersion = semver.New(0, 0, 0, "unknown", "")
	}
	config := updater.Config{
		Source:            github.MakeGithubProvider(viper.GetString("github-release-repo"), "checksums.txt"),
		CurrentVersion:    currentVersion,
		IncludePrerelease: currentVersion.Prerelease() != "", // Currently only update to a prerelease if the current version is a prerelease too
	}
	updateType := getUpdateType()
	// Some builds cannot (or should not) auto-update
	if updateType != nil {
		config.File = updateType.ArtifactName
		config.Apply = updateType.Apply
	}
	return config
}

var (
	updateCheckTicker *time.Ticker
	updateCheckStop   = make(chan bool)
)

func CheckForUpdate() error {
	if !shouldUseUpdater() {
		return nil
	}
	return Updater.CheckForUpdate()
}

func CheckInterval(interval time.Duration) {
	if !shouldUseUpdater() {
		return
	}
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
	if !shouldUseUpdater() {
		return nil
	}
	close(updateCheckStop)
	if Updater == nil {
		// No updater for this build type
		return nil
	}
	return Updater.OnExit(restart)
}
