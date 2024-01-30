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

func Init(config Config) {
	updateConfig := updater.Config{
		Source:                   github.MakeGithubProvider(viper.GetString("github-release-repo")),
		CurrentVersion:           viper.GetString("version"),
		UpdateFoundCallback:      config.UpdateFoundCallback,
		DownloadProgressCallback: config.DownloadProgressCallback,
		UpdateReadyCallback:      config.UpdateReadyCallback,
	}
	updateType := getUpdateType()
	// Some builds cannot (or should not) auto-update
	if updateType != nil {
		updateConfig.File = updateType.ArtifactName
		updateConfig.Apply = updateType.Apply
	}
	Updater = updater.MakeUpdater(updateConfig)
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
	updateCheckTicker := time.NewTicker(interval)
	go func() {
		err := Updater.CheckForUpdate()
		if err != nil {
			slog.Error("failed to check for update", slog.Any("error", err))
		}
		for range updateCheckTicker.C {
			err := Updater.CheckForUpdate()
			if err != nil {
				slog.Error("failed to check for update", slog.Any("error", err))
			}
		}
	}()
}

func OnExit(restart bool) error {
	if Updater == nil {
		// No updater for this build type
		return nil
	}
	if Updater.PendingUpdate == nil {
		if restart {
			slog.Warn("restart requested but no update is present. exiting anyway")
		}
		return nil
	}

	// We do have an update, but it might not be ready yet
	if !Updater.PendingUpdate.Ready {
		// TODO: I'd like to use an event here, but the autoupdater currently only takes a callback in the config
		recheckTicker := time.NewTicker(1 * time.Second)
		for range recheckTicker.C {
			if Updater.PendingUpdate.Ready {
				break
			}
		}
	}
	return Updater.OnExit(restart)
}

type Config struct {
	UpdateFoundCallback      func(latestVersion string, changelogs map[string]string)
	DownloadProgressCallback func(bytesDownloaded, totalBytes int64)
	UpdateReadyCallback      func()
}
