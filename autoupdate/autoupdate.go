package autoupdate

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/SatisfactoryModManager/autoupdate/updater"
	"github.com/satisfactorymodding/SatisfactoryModManager/autoupdate/updater/apply"
	"github.com/satisfactorymodding/SatisfactoryModManager/autoupdate/updater/source/github"
	"github.com/spf13/viper"
)

var Updater *updater.Updater
var updateApply apply.Apply
var releaseFile string

func Init(config AutoUpdateConfig) {
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

func CheckInterval(interval time.Duration) {
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

type AutoUpdateConfig struct {
	UpdateFoundCallback      func(latestVersion string, changelogs map[string]string)
	DownloadProgressCallback func(bytesDownloaded, totalBytes int64)
	UpdateReadyCallback      func()
}
