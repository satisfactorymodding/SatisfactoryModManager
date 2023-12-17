package ficsitcli

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/satisfactorymodding/ficsit-cli/utils"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders"
)

func (f *FicsitCLI) validateInstall(installation *InstallationInfo, progressItem string) error {
	f.EmitModsChange()
	defer f.EmitModsChange()

	installChannel := make(chan cli.InstallUpdate)

	defer f.setProgress(f.progress)

	type modProgress struct {
		downloadProgress utils.GenericProgress
		extractProgress  utils.GenericProgress
		downloading      bool
		complete         bool
	}
	modProgresses := map[string]modProgress{}
	modProgressesMutex := &sync.Mutex{}

	progressTicker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-progressTicker.C:
				modProgressesMutex.Lock()

				totalDownload := int64(0)
				completedDownload := int64(0)
				totalExtracted := int64(0)
				completedExtracted := int64(0)
				downloadsInProgressCount := 0
				extractionsInProgressCount := 0

				for _, modProgress := range modProgresses {
					if !modProgress.complete {
						if modProgress.downloading {
							downloadsInProgressCount++
						} else {
							extractionsInProgressCount++
						}
					}
					completedDownload += modProgress.downloadProgress.Completed
					totalDownload += modProgress.downloadProgress.Total
					completedExtracted += modProgress.extractProgress.Completed
					totalExtracted += modProgress.extractProgress.Total
				}

				if downloadsInProgressCount > 0 {
					if totalDownload != 0 {
						f.setProgress(&Progress{
							Item:     progressItem,
							Message:  fmt.Sprintf("Downloading %d mods", downloadsInProgressCount),
							Progress: float64(completedDownload) / float64(totalDownload),
						})
					}
				} else if extractionsInProgressCount > 0 {
					if totalExtracted != 0 {
						f.setProgress(&Progress{
							Item:     progressItem,
							Message:  fmt.Sprintf("Extracting %d mods", extractionsInProgressCount),
							Progress: float64(completedExtracted) / float64(totalExtracted),
						})
					}
				}
				modProgressesMutex.Unlock()
			}
		}
	}()

	go func() {
		for update := range installChannel {
			modProgressesMutex.Lock()
			switch update.Type {
			case cli.InstallUpdateTypeModDownload:
				modProgresses[update.Item.Mod] = modProgress{
					downloadProgress: update.Progress,
					extractProgress:  modProgresses[update.Item.Mod].extractProgress,
					downloading:      true,
					complete:         false,
				}
			case cli.InstallUpdateTypeModExtract:
				modProgresses[update.Item.Mod] = modProgress{
					downloadProgress: modProgresses[update.Item.Mod].downloadProgress,
					extractProgress:  update.Progress,
					downloading:      false,
					complete:         false,
				}
			case cli.InstallUpdateTypeModComplete:
				modProgresses[update.Item.Mod] = modProgress{
					downloadProgress: modProgresses[update.Item.Mod].downloadProgress,
					extractProgress:  modProgresses[update.Item.Mod].extractProgress,
					downloading:      false,
					complete:         true,
				}
			}
			modProgressesMutex.Unlock()
		}
		progressTicker.Stop()
		close(done)
	}()

	installErr := installation.Installation.Install(f.ficsitCli, installChannel)
	if installErr != nil {
		return errors.Wrap(installErr, "Failed to install")
	}
	return nil
}

func (f *FicsitCLI) EmitModsChange() {
	lockfileMods, err := f.GetSelectedInstallLockfileMods()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load lockfile")
		return
	}
	wailsRuntime.EventsEmit(f.ctx, "lockfileMods", lockfileMods)
	wailsRuntime.EventsEmit(f.ctx, "manifestMods", f.GetSelectedInstallProfileMods())
}

func (f *FicsitCLI) EmitGlobals() {
	installInfos := make([]*installfinders.Installation, 0, len(f.installations))
	for _, install := range f.installations {
		installInfos = append(installInfos, install.Info)
	}
	wailsRuntime.EventsEmit(f.ctx, "installations", installInfos)
	profileNames := make([]string, 0, len(f.ficsitCli.Profiles.Profiles))
	for k := range f.ficsitCli.Profiles.Profiles {
		profileNames = append(profileNames, k)
	}
	wailsRuntime.EventsEmit(f.ctx, "profiles", profileNames)

	if f.selectedInstallation == nil {
		return
	}

	wailsRuntime.EventsEmit(f.ctx, "selectedInstallation", f.selectedInstallation.Installation.Path)
	wailsRuntime.EventsEmit(f.ctx, "selectedProfile", f.selectedInstallation.Installation.Profile)
	wailsRuntime.EventsEmit(f.ctx, "modsEnabled", !f.selectedInstallation.Installation.Vanilla)
}

func (f *FicsitCLI) InstallMod(mod string) error {
	if f.progress != nil {
		return errors.New("Another operation in progress")
	}

	if f.selectedInstallation == nil {
		return errors.New("No installation selected")
	}

	profileName := f.selectedInstallation.Installation.Profile
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

	installErr := f.validateInstall(f.selectedInstallation, mod)

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

	if f.selectedInstallation == nil {
		return errors.New("No installation selected")
	}

	profileName := f.selectedInstallation.Installation.Profile
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

	installErr := f.validateInstall(f.selectedInstallation, mod)

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

	if f.selectedInstallation == nil {
		return errors.New("No installation selected")
	}

	profileName := f.selectedInstallation.Installation.Profile
	profile := f.GetProfile(profileName)

	profile.RemoveMod(mod)

	f.progress = &Progress{
		Item:     mod,
		Message:  "Checking for mods that are no longer needed",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(f.selectedInstallation, mod)

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

	if f.selectedInstallation == nil {
		return errors.New("No installation selected")
	}

	profileName := f.selectedInstallation.Installation.Profile
	profile := f.GetProfile(profileName)

	profile.SetModEnabled(mod, true)

	f.progress = &Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(f.selectedInstallation, mod)

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

	if f.selectedInstallation == nil {
		return errors.New("No installation selected")
	}

	profileName := f.selectedInstallation.Installation.Profile
	profile := f.GetProfile(profileName)

	profile.SetModEnabled(mod, false)

	f.progress = &Progress{
		Item:     mod,
		Message:  "Checking for mods that are no longer needed",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(f.selectedInstallation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to disable mod: %s", mod)
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}
