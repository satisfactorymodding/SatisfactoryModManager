package autoupdate

import (
	"log/slog"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/app"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/source/github"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater"
)

type autoUpdate struct {
	Updater *updater.Updater

	enabled bool

	restart bool

	updateCheckTicker *time.Ticker
	updateCheckStop   chan bool
}

var Updater *autoUpdate

func Init() {
	if Updater != nil {
		return
	}
	Updater = &autoUpdate{
		Updater: updater.MakeUpdater(makeUpdaterConfig()),
		enabled: shouldUseUpdater(),
	}
	Updater.Updater.UpdateFound.On(func(update updater.PendingUpdate) {
		wailsRuntime.EventsEmit(app.Context, "updateAvailable", update.Version.String(), update.Changelogs)
	})
	Updater.Updater.DownloadProgress.On(func(progress updater.UpdateDownloadProgress) {
		wailsRuntime.EventsEmit(app.Context, "updateDownloadProgress", progress.BytesDownloaded, progress.BytesTotal)
	})
	Updater.Updater.UpdateReady.On(func(interface{}) {
		wailsRuntime.EventsEmit(app.Context, "updateReady")
	})
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

func (u *autoUpdate) CheckForUpdates() {
	if !u.enabled {
		return
	}
	err := u.Updater.CheckForUpdate()
	if err != nil {
		slog.Warn("failed to check for updates", slog.Any("error", err))
	}
}

func (u *autoUpdate) UpdateAndRestart() {
	if !u.enabled {
		return
	}
	u.restart = true
	wailsRuntime.Quit(app.Context)
}

func (u *autoUpdate) CheckForUpdate() error {
	if !u.enabled {
		return nil
	}
	return u.Updater.CheckForUpdate()
}

func (u *autoUpdate) CheckInterval(interval time.Duration) {
	if !u.enabled {
		return
	}
	if u.updateCheckTicker != nil {
		return
	}
	u.updateCheckTicker = time.NewTicker(interval)
	u.updateCheckStop = make(chan bool)
	go func() {
		err := Updater.CheckForUpdate()
		if err != nil {
			slog.Error("failed to check for update", slog.Any("error", err))
		}
		for range u.updateCheckTicker.C {
			select {
			case <-u.updateCheckStop:
				return
			case <-u.updateCheckTicker.C:
				err := Updater.CheckForUpdate()
				if err != nil {
					slog.Error("failed to check for update", slog.Any("error", err))
				}
			}
		}
	}()
}

func (u *autoUpdate) OnExit() error {
	if !u.enabled {
		return nil
	}
	if u.updateCheckTicker != nil {
		close(u.updateCheckStop)
	}
	if Updater == nil {
		// No updater for this build type
		return nil
	}
	return u.Updater.OnExit(u.restart)
}
