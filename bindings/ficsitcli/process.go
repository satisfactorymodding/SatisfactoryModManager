package ficsitcli

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (f *FicsitCLI) validateInstall(installation *InstallationInfo, progressItem string) error {
	defer f.EmitModsChange()

	installChannel := make(chan cli.InstallUpdate)
	defer close(installChannel)

	defer f.setProgress(f.progress)

	go func() {
		for data := range installChannel {
			if data.DownloadProgress < 1 {
				f.setProgress(&Progress{
					Item:     progressItem,
					Message:  "Downloading " + data.ModName,
					Progress: data.DownloadProgress,
				})
			} else {
				f.setProgress(&Progress{
					Item:     progressItem,
					Message:  "Extracting " + data.ModName,
					Progress: data.ExtractProgress,
				})
			}
		}
	}()

	_, resolveErr := installation.Installation.ResolveProfile(f.ficsitCli)
	if resolveErr != nil {
		return errors.Wrap(resolveErr, "Failed to resolve profile")
	}
	installErr := installation.Installation.Install(f.ficsitCli, installChannel)
	if installErr != nil {
		return errors.Wrap(installErr, "Failed to install")
	}
	return nil
}

func (f *FicsitCLI) EmitModsChange() {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)
	lockfile, err := installation.Installation.LockFile(f.ficsitCli)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load lockfile")
		return
	}
	wailsRuntime.EventsEmit(f.ctx, "lockfileMods", lockfile)
	wailsRuntime.EventsEmit(f.ctx, "manifestMods", profile.Mods)
	wailsRuntime.EventsEmit(f.ctx, "selectedProfile", profileName)
}

func (f *FicsitCLI) InstallMod(mod string) error {
	if f.progress != nil {
		return errors.New("Another operation in progress")
	}

	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profileErr := profile.AddMod(mod, ">=0.0.0")
	if profileErr != nil {
		return errors.Wrapf(profileErr, "Failed to add mod: %s@latest", mod)
	}

	f.progress = &Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to install mod: %s@latest", mod)
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) InstallModVersion(mod string, version string) error {
	if f.progress != nil {
		return errors.New("Another operation in progress")
	}

	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profileErr := profile.AddMod(mod, version)
	if profileErr != nil {
		return errors.Wrapf(profileErr, "Failed to add mod: %s@%s", mod, version)
	}

	f.progress = &Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to install mod: %s@%s", mod, version)
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) RemoveMod(mod string) error {
	if f.progress != nil {
		return errors.New("Another operation in progress")
	}

	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profile.RemoveMod(mod)

	f.progress = &Progress{
		Item:     mod,
		Message:  "Checking for mods that are no longer needed",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to remove mod: %s", mod)
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) EnableMod(mod string) error {
	if f.progress != nil {
		return errors.New("Another operation in progress")
	}

	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profile.SetModEnabled(mod, true)

	f.progress = &Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to enable mod: %s", mod)
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) DisableMod(mod string) error {
	if f.progress != nil {
		return errors.New("Another operation in progress")
	}

	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profile.SetModEnabled(mod, false)

	f.progress = &Progress{
		Item:     mod,
		Message:  "Checking for mods that are no longer needed",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to disable mod: %s", mod)
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}
