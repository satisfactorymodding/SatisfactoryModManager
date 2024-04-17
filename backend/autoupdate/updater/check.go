package updater

import (
	"fmt"
	"log/slog"

	"github.com/Masterminds/semver/v3"
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

	latestVersion, err := u.config.Source.GetLatestVersion(u.config.IncludePrerelease)
	if err != nil {
		return fmt.Errorf("failed to get latest version: %w", err)
	}

	latestSemver, err := semver.NewVersion(latestVersion)
	if err != nil {
		return fmt.Errorf("failed to parse latest version %s: %w", latestVersion, err)
	}

	if u.PendingUpdate != nil && u.PendingUpdate.Version != nil && u.PendingUpdate.Ready {
		if !latestSemver.GreaterThan(u.PendingUpdate.Version) {
			return nil
		}
	} else {
		if !latestSemver.GreaterThan(u.config.CurrentVersion) {
			return nil
		}
	}

	changelogs, err := u.config.Source.GetChangelogs()
	if err != nil {
		return fmt.Errorf("failed to get changelogs: %w", err)
	}

	newChangelogs := make(map[string]string)
	for version, changelog := range changelogs {
		changelogSemver, err := semver.NewVersion(version)
		if err != nil {
			return fmt.Errorf("failed to parse version: %w", err)
		}
		if changelogSemver.GreaterThan(u.config.CurrentVersion) && changelogSemver.Compare(latestSemver) <= 0 {
			newChangelogs[version] = changelog
		}
	}

	u.PendingUpdate = &PendingUpdate{
		Version:    latestSemver,
		Changelogs: newChangelogs,
		Ready:      false,
	}
	u.UpdateFound.Dispatch(*u.PendingUpdate)

	if u.config.File == "" || u.config.Apply == nil {
		slog.Debug("no update file or apply method specified, not downloading update")
		return nil
	}

	file, length, err := u.config.Source.GetFile(latestVersion, u.config.File)
	if err != nil {
		return fmt.Errorf("failed to get file %s of version %s: %w", u.config.File, latestVersion, err)
	}
	defer file.Close()

	var checksum []byte
	if u.config.Checksum != nil {
		checksum, err = u.config.Checksum.GetChecksumForFile(u.config.Source, latestVersion, u.config.File)
		if err != nil {
			return fmt.Errorf("failed to get checksum for file %s of version %s: %w", u.config.File, latestVersion, err)
		}
	}

	progress := func(bytesDownloaded, bytesTotal int64) {
		u.DownloadProgress.Dispatch(UpdateDownloadProgress{
			BytesDownloaded: bytesDownloaded,
			BytesTotal:      bytesTotal,
		})
	}
	p := &progressReader{Reader: file, progressCallback: progress, contentLength: length}

	err = u.config.Apply.Apply(p, checksum)
	if err != nil {
		return fmt.Errorf("failed to apply update: %w", err)
	}
	u.PendingUpdate.Ready = true
	u.UpdateReady.Dispatch(nil)
	return nil
}
