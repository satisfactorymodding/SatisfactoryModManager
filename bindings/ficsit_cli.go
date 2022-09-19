package bindings

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/ficsit-cli/cli"

	"github.com/satisfactorymodding/SatisfactoryModManager/install_finders"
)

type FicsitCLI struct {
	ctx                  context.Context
	ficsitCli            *cli.GlobalContext
	installations        []*InstallationInfo
	installFindErrors    []error
	selectedInstallation *InstallationInfo
	progress             *Progress
}

type InstallationInfo struct {
	Installation *cli.Installation             `json:"installation"`
	Info         *install_finders.Installation `json:"info"`
}

type Progress struct {
	Item     string  `json:"item"`
	Message  string  `json:"message"`
	Progress float64 `json:"progress"`
}

func MakeFicsitCLI() (*FicsitCLI, error) {
	f := &FicsitCLI{}

	ficsitCli, err := cli.InitCLI()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize ficsit-cli")
	}
	f.ficsitCli = ficsitCli

	return f, nil
}

func (f *FicsitCLI) startup(ctx context.Context) {
	f.ctx = ctx

	f.initInstallations()
}

func (f *FicsitCLI) initInstallations() {
	installs, findErrors := install_finders.FindInstallations()

	f.installFindErrors = findErrors
	f.installations = []*InstallationInfo{}
	f.ficsitCli.Installations.Installations = []*cli.Installation{}

	for _, install := range installs {
		ficsitCliInstall, err := f.ficsitCli.Installations.AddInstallation(f.ficsitCli, install.Path, "Default")
		if err != nil {
			wailsRuntime.LogErrorf(f.ctx, "Failed to add installation: %v", err)
			continue
		}
		f.installations = append(f.installations, &InstallationInfo{
			Installation: ficsitCliInstall,
			Info:         install,
		})
	}

	sort.Slice(f.installations, func(i, j int) bool {
		if f.installations[i].Info.Launcher != f.installations[j].Info.Launcher {
			return f.installations[i].Info.Launcher < f.installations[j].Info.Launcher
		}
		return f.installations[i].Info.Branch < f.installations[j].Info.Branch
	})

	f.selectedInstallation = f.installations[0]
}

func (f *FicsitCLI) GetInstallationsInfo() []*InstallationInfo {
	return f.installations
}

func (f *FicsitCLI) GetInvalidInstalls() []string {
	var result []string
	for _, err := range f.installFindErrors {
		if casted, ok := err.(install_finders.InstallFindError); ok {
			result = append(result, casted.Path)
		}
	}
	return result
}

func (f *FicsitCLI) GetInstallation(path string) *InstallationInfo {
	for _, install := range f.installations {
		if install.Info.Path == path {
			return install
		}
	}

	return nil
}

func (f *FicsitCLI) SelectInstall(path string) {
	f.selectedInstallation = f.GetInstallation(path)
	f.emitModsChange()
}

func (f *FicsitCLI) GetSelectedInstall() *InstallationInfo {
	return f.selectedInstallation
}

func (f *FicsitCLI) SetProfile(profile string) {
	f.GetInstallation(f.selectedInstallation.Info.Path).Installation.SetProfile(f.ficsitCli, profile)
	f.emitModsChange()
}

func (f *FicsitCLI) GetSelectedProfile() *string {
	if f.selectedInstallation == nil {
		return nil
	}
	return &f.selectedInstallation.Installation.Profile
}

func (f *FicsitCLI) GetProfiles() []string {
	profileNames := make([]string, 0, len(f.ficsitCli.Profiles.Profiles))
	for k := range f.ficsitCli.Profiles.Profiles {
		profileNames = append(profileNames, k)
	}
	sort.Strings(profileNames)
	return profileNames
}

func (f *FicsitCLI) GetProfile(profile string) *cli.Profile {
	return f.ficsitCli.Profiles.GetProfile(profile)
}

func (f *FicsitCLI) validateInstall(installation *InstallationInfo, progressItem string) error {
	installChannel := make(chan cli.InstallUpdate)

	defer f.SetProgress(f.progress)

	go func() {
		for data := range installChannel {
			if data.DownloadProgress < 1 {
				f.SetProgress(&Progress{
					Item:     progressItem,
					Message:  "Downloading " + data.ModName,
					Progress: data.DownloadProgress,
				})
			} else {
				f.SetProgress(&Progress{
					Item:     progressItem,
					Message:  "Extracting " + data.ModName,
					Progress: data.DownloadProgress,
				})
			}
		}
	}()

	installErr := installation.Installation.Install(f.ficsitCli, installChannel)

	close(installChannel)

	return installErr
}

func (f *FicsitCLI) InstallMod(mod string) error {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profileErr := profile.AddMod(mod, ">=0.0.0")
	if profileErr != nil {
		return errors.Wrapf(profileErr, "Failed to add mod: %s@latest", mod)
	}

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	})

	defer f.SetProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to install mod: %s@latest", mod)
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()

	return nil
}

func (f *FicsitCLI) InstallModVersion(mod string, version string) error {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profileErr := profile.AddMod(mod, version)
	if profileErr != nil {
		return errors.Wrapf(profileErr, "Failed to add mod: %s@%s", mod, version)
	}

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	})

	defer f.SetProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to install mod: %s@%s", mod, version)
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()

	return nil
}

func (f *FicsitCLI) RemoveMod(mod string) error {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profile.RemoveMod(mod)

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Checking for mods that are no longer needed",
		Progress: -1,
	})

	defer f.SetProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to remove mod: %s", mod)
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()

	return nil
}

func (f *FicsitCLI) EnableMod(mod string) error {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profile.SetModEnabled(mod, true)

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	})

	defer f.SetProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to enable mod: %s", mod)
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()

	return nil
}

func (f *FicsitCLI) DisableMod(mod string) error {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profile.SetModEnabled(mod, false)

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Checking for mods that are no longer needed",
		Progress: -1,
	})

	defer f.SetProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		return errors.Wrapf(installErr, "Failed to disable mod: %s", mod)
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()

	return nil
}

func (f *FicsitCLI) AddProfile(name string) error {
	_, err := f.ficsitCli.Profiles.AddProfile(name)

	if err != nil {
		return errors.Wrapf(err, "Failed to add profile: %s", name)
	}

	f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) RenameProfile(oldName string, newName string) error {
	err := f.ficsitCli.Profiles.RenameProfile(f.ficsitCli, oldName, newName)

	if err != nil {
		return errors.Wrapf(err, "Failed to rename profile: %s -> %s", oldName, newName)
	}

	f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) DeleteProfile(name string) error {
	err := f.ficsitCli.Profiles.DeleteProfile(name)

	if err != nil {
		return errors.Wrapf(err, "Failed to delete profile: %s", name)
	}

	f.ficsitCli.Profiles.Save()

	for _, install := range f.installations {
		if install.Installation.Profile == name {
			install.Installation.Profile = cli.DefaultProfileName
		}
	}

	return nil
}

type Update struct {
	Item           string `json:"item"`
	CurrentVersion string `json:"currentVersion"`
	NewVersion     string `json:"newVersion"`
}

func (f *FicsitCLI) CheckForUpdates() ([]Update, error) {
	currentLockfile, err := f.selectedInstallation.Installation.LockFile(f.ficsitCli)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get current lockfile")
	}

	if currentLockfile == nil {
		return nil, nil
	}

	profile := f.GetProfile(f.selectedInstallation.Installation.Profile)

	resolver := cli.NewDependencyResolver(f.ficsitCli.APIClient)

	gameVersion, err := f.selectedInstallation.Installation.GetGameVersion(f.ficsitCli)
	if err != nil {
		return nil, errors.Wrap(err, "failed to detect game version")
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
		return nil, errors.Wrap(err, "error resolving dependencies")
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

	previousLockfile, err := f.selectedInstallation.Installation.LockFile(f.ficsitCli)
	if err != nil {
		return errors.Wrap(err, "Failed to get current lockfile")
	}

	profile := f.GetProfile(f.selectedInstallation.Installation.Profile)
	for modReference, modData := range profile.Mods {
		profile.Mods[modReference] = cli.ProfileMod{
			Enabled: modData.Enabled,
			Version: ">=0.0.0",
		}
	}

	f.selectedInstallation.Installation.WriteLockFile(f.ficsitCli, cli.LockFile{})

	f.SetProgress(&Progress{
		Item:     "__update__",
		Message:  "Updating...",
		Progress: -1,
	})

	defer f.SetProgress(nil)

	err = f.validateInstall(f.selectedInstallation, "__update__")

	if err != nil {
		f.selectedInstallation.Installation.WriteLockFile(f.ficsitCli, *previousLockfile)
		return errors.Wrap(err, "Failed to update mods")
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()

	return nil
}

func (f *FicsitCLI) GetCurrentLockfile(install *InstallationInfo) (*cli.LockFile, error) {
	lockfile, err := install.Installation.LockFile(f.ficsitCli)
	if err != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed read lockfile: %v", err)
		return nil, err
	}
	return lockfile, nil
}

func (f *FicsitCLI) emitModsChange() {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)
	lockfile, err := installation.Installation.LockFile(f.ficsitCli)
	if err != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed read lockfile: %v", err)
		return
	}
	wailsRuntime.EventsEmit(f.ctx, "lockfileMods", lockfile)
	wailsRuntime.EventsEmit(f.ctx, "manifestMods", profile.Mods)
}

func (f *FicsitCLI) SetProgress(progress *Progress) {
	f.progress = progress
	wailsRuntime.EventsEmit(f.ctx, "progress", f.progress)
}
