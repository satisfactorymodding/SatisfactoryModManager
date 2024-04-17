package updater

import (
	"fmt"
	"sync"

	"github.com/Masterminds/semver/v3"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

type UpdateDownloadProgress struct {
	BytesDownloaded, BytesTotal int64
}

type Updater struct {
	config        Config
	lock          sync.Mutex
	PendingUpdate *PendingUpdate

	UpdateFound      utils.EventDispatcher[PendingUpdate]
	DownloadProgress utils.EventDispatcher[UpdateDownloadProgress]
	UpdateReady      utils.EventDispatcher[interface{}]
}

type PendingUpdate struct {
	Version    *semver.Version
	Changelogs map[string]string
	Ready      bool
}

type Config struct {
	Source            Source
	File              string
	Checksum          ChecksumSource
	Apply             Apply
	CurrentVersion    *semver.Version
	IncludePrerelease bool
}

func MakeUpdater(config Config) *Updater {
	return &Updater{
		config: config,
	}
}

func (u *Updater) OnExit(restart bool) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.PendingUpdate == nil {
		if restart {
			return fmt.Errorf("restart requested but no update is present")
		}
		return nil
	}

	// We do have an update
	// and since we have the lock, we can be sure that no other update can be found while we're here

	// Though, applying the update might have errored, meaning the update is not actually ready
	if !u.PendingUpdate.Ready {
		if restart {
			return fmt.Errorf("restart requested but update is not ready")
		}
		return nil
	}

	// Now the update is definitely ready
	return u.config.Apply.OnExit(restart)
}
