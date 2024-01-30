//go:build !bindings

package updater

import (
	"log/slog"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

func (u *Updater) CheckForUpdate() error {
	if u.PendingUpdate != nil {
		u.UpdateFound.Dispatch(*u.PendingUpdate)
		if u.PendingUpdate.Ready {
			u.UpdateReady.Dispatch(nil)
		}
	}

	u.lock.Lock()
	defer u.lock.Unlock()

	latestVersion, err := u.config.Source.GetLatestVersion()
	if err != nil {
		return errors.Wrap(err, "failed to get latest version")
	}

	latestSemver, err := semver.NewVersion(latestVersion)
	if err != nil {
		return errors.Wrapf(err, "failed to parse latest version %s", latestVersion)
	}

	var pendingSemver *semver.Version
	if u.PendingUpdate != nil {
		pendingSemver, err = semver.NewVersion(u.PendingUpdate.Version)
		if err != nil {
			return errors.Wrapf(err, "failed to parse pending version %s", u.PendingUpdate.Version)
		}
	}
	currentSemver, err := semver.NewVersion(u.config.CurrentVersion)
	if err != nil {
		return errors.Wrapf(err, "failed to parse current version %s", u.config.CurrentVersion)
	}

	if pendingSemver != nil {
		if !latestSemver.GreaterThan(pendingSemver) {
			return nil
		}
	} else {
		if !latestSemver.GreaterThan(currentSemver) {
			return nil
		}
	}

	changelogs, err := u.config.Source.GetChangelogs()
	if err != nil {
		return errors.Wrap(err, "failed to get changelogs")
	}

	newChangelogs := make(map[string]string)
	for version, changelog := range changelogs {
		semver, err := semver.NewVersion(version)
		if err != nil {
			return errors.Wrap(err, "failed to parse version")
		}
		if semver.GreaterThan(currentSemver) && semver.Compare(latestSemver) <= 0 {
			newChangelogs[version] = changelog
		}
	}

	u.PendingUpdate = &PendingUpdate{
		Version:    latestVersion,
		Changelogs: newChangelogs,
		Ready:      false,
	}
	u.UpdateFound.Dispatch(*u.PendingUpdate)

	if u.config.File == "" || u.config.Apply == nil {
		slog.Debug("no update file or apply method specified, not downloading update")
		return nil
	}

	file, length, err := u.config.Source.GetFile(u.config.File)
	if err != nil {
		return errors.Wrap(err, "failed to get file")
	}
	defer file.Close()

	progress := func(bytesDownloaded, bytesTotal int64) {
		u.DownloadProgress.Dispatch(UpdateDownloadProgress{
			BytesDownloaded: bytesDownloaded,
			BytesTotal:      bytesTotal,
		})
	}
	p := &progressReader{Reader: file, progressCallback: progress, contentLength: length}

	err = u.config.Apply.Apply(p)
	if err != nil {
		return errors.Wrap(err, "failed to apply update")
	}
	u.PendingUpdate.Ready = true
	u.UpdateReady.Dispatch(nil)
	return nil
}
