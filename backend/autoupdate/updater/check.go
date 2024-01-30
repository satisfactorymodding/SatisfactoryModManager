//go:build !bindings

package updater

import (
	"log/slog"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

func (u *Updater) CheckForUpdate() error {
	if u.PendingUpdate != nil {
		u.config.UpdateFoundCallback(u.PendingUpdate.Version, u.PendingUpdate.Changelogs)
		if u.PendingUpdate.Ready {
			u.config.UpdateReadyCallback()
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
		if semver.GreaterThan(currentSemver) {
			newChangelogs[version] = changelog
		}
	}

	u.PendingUpdate = &PendingUpdate{
		Version:    latestVersion,
		Changelogs: newChangelogs,
		Ready:      false,
	}
	u.config.UpdateFoundCallback(latestVersion, newChangelogs)

	if u.config.File == "" || u.config.Apply == nil {
		slog.Debug("no update file or apply method specified, not downloading update")
		return nil
	}

	file, length, err := u.config.Source.GetFile(u.config.File)
	if err != nil {
		return errors.Wrap(err, "failed to get file")
	}
	defer file.Close()
	p := &progressReader{Reader: file, progressCallback: u.config.DownloadProgressCallback, contentLength: length}

	err = u.config.Apply.Apply(p)
	if err != nil {
		return errors.Wrap(err, "failed to apply update")
	}
	u.PendingUpdate.Ready = true
	u.config.UpdateReadyCallback()
	return nil
}
