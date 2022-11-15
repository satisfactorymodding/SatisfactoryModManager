package updater

import (
	"sync"

	"github.com/satisfactorymodding/SatisfactoryModManager/autoupdate/updater/apply"
	"github.com/satisfactorymodding/SatisfactoryModManager/autoupdate/updater/source"
)

type Updater struct {
	config        Config
	lock          sync.Mutex
	PendingUpdate *PendingUpdate
}

type PendingUpdate struct {
	Version    string
	Changelogs map[string]string
	Ready      bool
}

type Config struct {
	Source                   source.Source
	Apply                    apply.Apply
	CurrentVersion           string
	UpdateFoundCallback      func(latestVersion string, changelogs map[string]string)
	DownloadProgressCallback func(bytesDownloaded, totalBytes int64)
	UpdateReadyCallback      func()
}

func MakeUpdater(config Config) *Updater {
	return &Updater{
		config: config,
	}
}
