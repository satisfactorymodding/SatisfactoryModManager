package updater

import (
	"sync"
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
	Source                   Source
	File                     string
	Apply                    Apply
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

func (u *Updater) OnExit(restart bool) error {
	return u.config.Apply.OnExit(restart)
}
