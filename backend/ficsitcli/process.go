package ficsitcli

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	ficsitUtils "github.com/satisfactorymodding/ficsit-cli/utils"
	resolver "github.com/satisfactorymodding/ficsit-resolver"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	appCommon "github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
)

func (f *ficsitCLI) validateInstall(installation *cli.Installation) error {
	if !f.isValidInstall(installation.Path) {
		return fmt.Errorf("invalid installation: %s", installation.Path)
	}

	f.EmitModsChange()
	defer f.EmitModsChange()

	installChannel := make(chan cli.InstallUpdate)

	defer f.setProgress(f.progress)

	type modProgress struct {
		downloadProgress ficsitUtils.GenericProgress
		extractProgress  ficsitUtils.GenericProgress
		downloading      bool
		complete         bool
	}
	modProgresses := xsync.NewMapOf[string, modProgress]()

	progressTicker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)
	defer progressTicker.Stop()
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-progressTicker.C:
				newTasks := make(map[string]ProgressTask)

				modProgresses.Range(func(modAndVersion string, progress modProgress) bool {
					if progress.downloadProgress.Total != 0 {
						newTasks[fmt.Sprintf("%s:download", modAndVersion)] = ProgressTask{
							Current: progress.downloadProgress.Completed,
							Total:   progress.downloadProgress.Total,
						}
					}
					if progress.downloadProgress.Total != 0 || progress.extractProgress.Total != 0 {
						newTasks[fmt.Sprintf("%s:extract", modAndVersion)] = ProgressTask{
							Current: progress.extractProgress.Completed,
							Total:   progress.extractProgress.Total,
						}
					}
					return true
				})

				f.setProgress(&Progress{
					Action: f.progress.Action,
					Item:   f.progress.Item,
					Tasks:  newTasks,
				})
			}
		}
	}()

	go func() {
		for update := range installChannel {
			if update.Type == cli.InstallUpdateTypeOverall {
				// Although this wouldn't cause any issues in the progress generation above, we can ignore this update.
				continue
			}
			itemName := fmt.Sprintf("%s@%s", update.Item.Mod, update.Item.Version)
			modProgresses.Compute(itemName, func(oldValue modProgress, loaded bool) (modProgress, bool) {
				if oldValue.complete {
					// Sometimes extract updates are received after the mod is marked as complete.
					return oldValue, false
				}
				oldValue.complete = update.Type == cli.InstallUpdateTypeModComplete
				oldValue.downloading = update.Type == cli.InstallUpdateTypeModDownload

				switch update.Type {
				case cli.InstallUpdateTypeModDownload:
					oldValue.downloadProgress = update.Progress
				case cli.InstallUpdateTypeModExtract:
					oldValue.extractProgress = update.Progress
				}
				return oldValue, false
			})
		}
	}()

	installErr := installation.Install(f.ficsitCli, installChannel)
	if installErr != nil {
		var solvingError resolver.DependencyResolverError
		if errors.As(installErr, &solvingError) {
			return solvingError
		}
		return installErr //nolint:wrapcheck
	}
	return nil
}

func (f *ficsitCLI) EmitModsChange() {
	lockfileMods, err := f.GetSelectedInstallLockfileMods()
	if err != nil {
		slog.Error("failed to load lockfile", slog.Any("error", err))
		return
	}
	wailsRuntime.EventsEmit(appCommon.AppContext, "lockfileMods", lockfileMods)
	wailsRuntime.EventsEmit(appCommon.AppContext, "manifestMods", f.GetSelectedInstallProfileMods())
}

func (f *ficsitCLI) EmitGlobals() {
	if appCommon.AppContext == nil {
		// This function can be called from AddRemoteServer, which is used during initialization
		// at which point the context is not set yet.
		// We can safely ignore this call.
		return
	}
	wailsRuntime.EventsEmit(appCommon.AppContext, "installations", f.GetInstallations())
	wailsRuntime.EventsEmit(appCommon.AppContext, "installationsMetadata", f.GetInstallationsMetadata())
	wailsRuntime.EventsEmit(appCommon.AppContext, "remoteServers", f.GetRemoteInstallations())
	profileNames := make([]string, 0, len(f.ficsitCli.Profiles.Profiles))
	for k := range f.ficsitCli.Profiles.Profiles {
		profileNames = append(profileNames, k)
	}
	wailsRuntime.EventsEmit(appCommon.AppContext, "profiles", profileNames)

	selectedInstallation := f.GetSelectedInstall()

	if selectedInstallation == nil {
		return
	}

	wailsRuntime.EventsEmit(appCommon.AppContext, "selectedInstallation", selectedInstallation.Path)
	wailsRuntime.EventsEmit(appCommon.AppContext, "selectedProfile", selectedInstallation.Profile)
	wailsRuntime.EventsEmit(appCommon.AppContext, "modsEnabled", !selectedInstallation.Vanilla)
}

func (f *ficsitCLI) InstallMod(mod string) error {
	if f.progress != nil {
		return fmt.Errorf("another operation in progress")
	}

	selectedInstallation := f.GetSelectedInstall()

	if selectedInstallation == nil {
		return fmt.Errorf("no installation selected")
	}

	l := slog.With(
		slog.String("task", "installMod"),
		slog.String("mod", mod),
		slog.String("install", selectedInstallation.Path),
		slog.String("profile", selectedInstallation.Profile),
	)

	profileName := selectedInstallation.Profile
	profile := f.GetProfile(profileName)

	profileErr := profile.AddMod(mod, ">=0.0.0")
	if profileErr != nil {
		l.Error("failed to add mod", slog.Any("error", profileErr))
		return fmt.Errorf("failed to add mod: %s@latest: %w", mod, profileErr)
	}

	err := f.ficsitCli.Profiles.Save()
	if err != nil {
		l.Error("failed to save profile", slog.Any("error", err))
	}

	f.progress = newProgress(ActionInstall, newSimpleItem(mod))

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstallation)

	if installErr != nil {
		l.Error("failed to install", slog.Any("error", installErr))
		return installErr
	}

	return nil
}

func (f *ficsitCLI) InstallModVersion(mod string, version string) error {
	if f.progress != nil {
		return fmt.Errorf("another operation in progress")
	}

	selectedInstallation := f.GetSelectedInstall()

	if selectedInstallation == nil {
		return fmt.Errorf("no installation selected")
	}

	l := slog.With(
		slog.String("task", "installModVersion"),
		slog.String("mod", mod),
		slog.String("version", version),
		slog.String("install", selectedInstallation.Path),
		slog.String("profile", selectedInstallation.Profile),
	)

	profile := f.GetProfile(selectedInstallation.Profile)

	profileErr := profile.AddMod(mod, version)
	if profileErr != nil {
		l.Error("failed to add mod", slog.Any("error", profileErr))
		return fmt.Errorf("failed to add mod: %s@%s: %w", mod, version, profileErr)
	}

	err := f.ficsitCli.Profiles.Save()
	if err != nil {
		l.Error("failed to save profile", slog.Any("error", err))
	}

	f.progress = newProgress(ActionInstall, newItem(mod, version))

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstallation)

	if installErr != nil {
		l.Error("failed to install", slog.Any("error", installErr))
		return installErr
	}

	return nil
}

func (f *ficsitCLI) RemoveMod(mod string) error {
	if f.progress != nil {
		return fmt.Errorf("another operation in progress")
	}

	selectedInstallation := f.GetSelectedInstall()

	if selectedInstallation == nil {
		return fmt.Errorf("no installation selected")
	}

	l := slog.With(
		slog.String("task", "removeMod"),
		slog.String("mod", mod),
		slog.String("install", selectedInstallation.Path),
		slog.String("profile", selectedInstallation.Profile),
	)

	profile := f.GetProfile(selectedInstallation.Profile)

	profile.RemoveMod(mod)

	err := f.ficsitCli.Profiles.Save()
	if err != nil {
		l.Error("failed to save profile", slog.Any("error", err))
	}

	f.progress = newProgress(ActionUninstall, newSimpleItem(mod))

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstallation)

	if installErr != nil {
		l.Error("failed to install", slog.Any("error", installErr))
		return installErr
	}

	return nil
}

func (f *ficsitCLI) EnableMod(mod string) error {
	if f.progress != nil {
		return fmt.Errorf("another operation in progress")
	}

	selectedInstallation := f.GetSelectedInstall()

	if selectedInstallation == nil {
		return fmt.Errorf("no installation selected")
	}

	l := slog.With(
		slog.String("task", "enableMod"),
		slog.String("mod", mod),
		slog.String("install", selectedInstallation.Path),
		slog.String("profile", selectedInstallation.Profile),
	)

	profile := f.GetProfile(selectedInstallation.Profile)

	profile.SetModEnabled(mod, true)

	err := f.ficsitCli.Profiles.Save()
	if err != nil {
		l.Error("failed to save profile", slog.Any("error", err))
	}

	f.progress = newProgress(ActionEnable, newSimpleItem(mod))

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstallation)

	if installErr != nil {
		l.Error("failed to install", slog.Any("error", installErr))
		return installErr
	}

	return nil
}

func (f *ficsitCLI) DisableMod(mod string) error {
	if f.progress != nil {
		return fmt.Errorf("another operation in progress")
	}

	selectedInstallation := f.GetSelectedInstall()

	if selectedInstallation == nil {
		return fmt.Errorf("no installation selected")
	}

	l := slog.With(
		slog.String("task", "disableMod"),
		slog.String("mod", mod),
		slog.String("install", selectedInstallation.Path),
		slog.String("profile", selectedInstallation.Profile),
	)

	profile := f.GetProfile(selectedInstallation.Profile)

	profile.SetModEnabled(mod, false)

	err := f.ficsitCli.Profiles.Save()
	if err != nil {
		l.Error("failed to save profile", slog.Any("error", err))
	}

	f.progress = newProgress(ActionDisable, newSimpleItem(mod))

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstallation)

	if installErr != nil {
		l.Error("failed to install", slog.Any("error", installErr))
		return installErr
	}

	return nil
}
