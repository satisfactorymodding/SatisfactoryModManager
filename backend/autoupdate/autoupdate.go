package autoupdate

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater/apply"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater/source/github"
)

var (
	Updater     *updater.Updater
	updateApply apply.Apply
	releaseFile string
)

func Init(config Config) {
	releaseFile, updateApply = getInstallType()
	Updater = updater.MakeUpdater(updater.Config{
		Source:                   github.MakeGithubProvider(viper.GetString("github-release-repo"), releaseFile),
		Apply:                    updateApply,
		CurrentVersion:           viper.GetString("version"),
		UpdateFoundCallback:      config.UpdateFoundCallback,
		DownloadProgressCallback: config.DownloadProgressCallback,
		UpdateReadyCallback:      config.UpdateReadyCallback,
	})
}

var checkStarted bool

func CheckInterval(interval time.Duration) {
	if checkStarted {
		return
	}
	checkStarted = true
	updateCheckTicker := time.NewTicker(interval)
	go func() {
		err := Updater.CheckForUpdate()
		if err != nil {
			log.Error().Err(err).Msg("Failed to check for update")
		}
		for range updateCheckTicker.C {
			err := Updater.CheckForUpdate()
			if err != nil {
				log.Error().Err(err).Msg("Failed to check for update")
			}
		}
	}()
}

func OnExit(restart bool) error {
	return updateApply.OnExit(restart)
}

type Config struct {
	UpdateFoundCallback      func(latestVersion string, changelogs map[string]string)
	DownloadProgressCallback func(bytesDownloaded, totalBytes int64)
	UpdateReadyCallback      func()
}
