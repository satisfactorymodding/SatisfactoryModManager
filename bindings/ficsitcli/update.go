package ficsitcli

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/ficsit-cli/cli"
)

type Update struct {
	Item           string `json:"item"`
	CurrentVersion string `json:"currentVersion"`
	NewVersion     string `json:"newVersion"`
}

func (f *FicsitCLI) CheckForUpdates() ([]Update, error) {
	l := log.With().Str("task", "checkForUpdates").Logger()

	currentLockfile, err := f.selectedInstallation.Installation.LockFile(f.ficsitCli)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get current lockfile")
		return nil, errors.Wrap(err, "Failed to get current lockfile")
	}

	if currentLockfile == nil {
		return nil, nil
	}

	profile := f.GetProfile(f.selectedInstallation.Installation.Profile)

	resolver := cli.NewDependencyResolver(f.ficsitCli.Provider)

	gameVersion, err := f.selectedInstallation.Installation.GetGameVersion(f.ficsitCli)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get game version")
		return nil, errors.Wrap(err, "Failed to get game version")
	}

	updateProfile := &cli.Profile{
		Name: "Update temp",
		Mods: make(map[string]cli.ProfileMod),
	}
	for modReference, modData := range profile.Mods {
		updateProfile.Mods[modReference] = cli.ProfileMod{
			Enabled: modData.Enabled,
			Version: ">=0.0.0",
		}
	}
	newLockfile, err := updateProfile.Resolve(resolver, nil, gameVersion)
	if err != nil {
		l.Error().Err(err).Msg("Failed to resolve new lockfile")
		return nil, errors.Wrap(err, "Error resolving dependencies")
	}

	updates := []Update{}

	for modReference, newLockedMod := range newLockfile {
		if prevLockedMod, ok := (*currentLockfile)[modReference]; ok {
			if newLockedMod.Version != prevLockedMod.Version {
				updates = append(updates, Update{
					Item:           modReference,
					CurrentVersion: prevLockedMod.Version,
					NewVersion:     newLockedMod.Version,
				})
			}
		}
	}

	return updates, nil
}

func (f *FicsitCLI) UpdateAllMods() error {
	l := log.With().Str("task", "updateAllMods").Logger()

	if f.progress != nil {
		l.Error().Msg("Another operation in progress")
		return errors.New("Another operation in progress")
	}

	previousLockfile, err := f.selectedInstallation.Installation.LockFile(f.ficsitCli)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get current lockfile")
		return errors.Wrap(err, "Failed to get current lockfile")
	}

	profile := f.GetProfile(f.selectedInstallation.Installation.Profile)
	for modReference, modData := range profile.Mods {
		profile.Mods[modReference] = cli.ProfileMod{
			Enabled: modData.Enabled,
			Version: ">=0.0.0",
		}
	}

	_ = f.selectedInstallation.Installation.WriteLockFile(f.ficsitCli, cli.LockFile{})

	f.progress = &Progress{
		Item:     "__update__",
		Message:  "Updating...",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	err = f.validateInstall(f.selectedInstallation, "__update__")

	if err != nil {
		_ = f.selectedInstallation.Installation.WriteLockFile(f.ficsitCli, *previousLockfile)
		l.Error().Err(err).Msg("Failed to validate installation")
		return errors.Wrap(err, "Failed to update mods")
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}
