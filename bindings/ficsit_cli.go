package bindings

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/ficsit-cli/cli"

	"github.com/satisfactorymodding/SatisfactoryModManager/install_finders"
	"github.com/satisfactorymodding/SatisfactoryModManager/settings"
)

type FicsitCLI struct {
	ctx                  context.Context
	ficsitCli            *cli.GlobalContext
	installations        []*InstallationInfo
	installFindErrors    []error
	selectedInstallation *InstallationInfo
	progress             *Progress
	isGameRunning        bool
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

	ficsitCli, err := cli.InitCLI(false)
	if err != nil {
		err = errors.Wrap(err, "failed to initialize ficsit-cli")
		wailsRuntime.LogError(f.ctx, err.Error())
		return nil, err
	}
	f.ficsitCli = ficsitCli

	return f, nil
}

func (f *FicsitCLI) startup(ctx context.Context) {
	f.ctx = ctx

	f.initInstallations()

	gameRunningTIcker := time.NewTicker(5 * time.Second)
	go func() {
		for range gameRunningTIcker.C {
			processes, err := ps.Processes()
			if err != nil {
				wailsRuntime.LogErrorf(f.ctx, "Failed to get processes: %v", err)
				continue
			}
			f.isGameRunning = false
			for _, process := range processes {
				if process.Executable() == "FactoryGame-Win64-Shipping.exe" || process.Executable() == "FactoryGame-Win64-Shipping" {
					f.isGameRunning = true
					break
				}
			}
			wailsRuntime.EventsEmit(f.ctx, "isGameRunning", f.isGameRunning)
		}
	}()
}

func (f *FicsitCLI) LaunchGame() {
	if f.selectedInstallation == nil {
		wailsRuntime.LogErrorf(f.ctx, "No installation selected")
		return
	}
	cmd := exec.Command(f.selectedInstallation.Info.LaunchPath[0], f.selectedInstallation.Info.LaunchPath[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to launch game (%s): error: %v, output: %s", cmd, err, out)
		return
	}
	f.isGameRunning = true
	wailsRuntime.EventsEmit(f.ctx, "isGameRunning", f.isGameRunning)
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

	for _, install := range f.installations {
		savedSelectedProfile, ok := settings.Settings.SelectedProfile[install.Info.Path]
		if ok {
			err := install.Installation.SetProfile(f.ficsitCli, savedSelectedProfile)
			if err != nil {
				wailsRuntime.LogErrorf(f.ctx, "Failed to set saved profile: %v", err)
			}
		}
	}

	f.selectedInstallation = f.installations[0]

	savedSelectedInstall := settings.Settings.SelectedInstall
	if savedSelectedInstall != "" {
		for _, install := range f.installations {
			if install.Info.Path == savedSelectedInstall {
				f.selectedInstallation = install
				break
			}
		}
	}
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
	installation := f.GetInstallation(path)
	if installation == nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to find installation %s", path)
		return
	}
	f.selectedInstallation = installation
	f.emitModsChange()
	settings.Settings.SelectedInstall = installation.Info.Path
	settings.SaveSettings()
}

func (f *FicsitCLI) GetSelectedInstall() *InstallationInfo {
	return f.selectedInstallation
}

func (f *FicsitCLI) SetProfile(profile string) {
	if f.selectedInstallation == nil {
		wailsRuntime.LogErrorf(f.ctx, "No installation selected")
		return
	}
	err := f.selectedInstallation.Installation.SetProfile(f.ficsitCli, profile)
	if err != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to set profile: %v", err)
		return
	}
	f.emitModsChange()
	settings.Settings.SelectedProfile[f.selectedInstallation.Info.Path] = profile
	settings.SaveSettings()
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
		profileErr := errors.Wrapf(profileErr, "Failed to add mod: %s@latest", mod)
		wailsRuntime.LogError(f.ctx, profileErr.Error())
		return profileErr
	}

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	})

	defer f.SetProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		installErr := errors.Wrapf(installErr, "Failed to install mod: %s@latest", mod)
		wailsRuntime.LogError(f.ctx, installErr.Error())
		return installErr
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
		profileErr := errors.Wrapf(profileErr, "Failed to add mod: %s@%s", mod, version)
		wailsRuntime.LogError(f.ctx, profileErr.Error())
		return profileErr
	}

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	})

	defer f.SetProgress(nil)

	installErr := f.validateInstall(installation, mod)

	if installErr != nil {
		installErr := errors.Wrapf(installErr, "Failed to install mod: %s@%s", mod, version)
		wailsRuntime.LogError(f.ctx, installErr.Error())
		return installErr
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
		installErr := errors.Wrapf(installErr, "Failed to remove mod: %s", mod)
		wailsRuntime.LogError(f.ctx, installErr.Error())
		return installErr
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
		installErr := errors.Wrapf(installErr, "Failed to enable mod: %s", mod)
		wailsRuntime.LogError(f.ctx, installErr.Error())
		return installErr
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
		installErr := errors.Wrapf(installErr, "Failed to disable mod: %s", mod)
		wailsRuntime.LogError(f.ctx, installErr.Error())
		return installErr
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()

	return nil
}

func (f *FicsitCLI) AddProfile(name string) error {
	_, err := f.ficsitCli.Profiles.AddProfile(name)

	if err != nil {
		err := errors.Wrapf(err, "Failed to add profile: %s", name)
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) RenameProfile(oldName string, newName string) error {
	err := f.ficsitCli.Profiles.RenameProfile(f.ficsitCli, oldName, newName)

	if err != nil {
		err := errors.Wrapf(err, "Failed to rename profile: %s -> %s", oldName, newName)
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) DeleteProfile(name string) error {
	err := f.ficsitCli.Profiles.DeleteProfile(name)

	if err != nil {
		err := errors.Wrapf(err, "Failed to delete profile: %s", name)
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	f.ficsitCli.Profiles.Save()

	for _, install := range f.installations {
		if install.Installation.Profile == name {
			install.Installation.Profile = cli.DefaultProfileName
		}
	}

	return nil
}

type ExportedProfile struct {
	Profile  *cli.Profile             `json:"profile"`
	LockFile *cli.LockFile            `json:"lockfile"`
	Metadata *ExportedProfileMetadata `json:"metadata"`
}

type ExportedProfileMetadata struct {
	GameVersion int `json:"gameVersion"`
}

func (f *FicsitCLI) ExportCurrentProfile() error {
	selectedInstall := f.GetSelectedInstall()
	if selectedInstall == nil {
		err := errors.New("No installation selected")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	profileName := f.GetSelectedProfile()
	if profileName == nil {
		err := errors.New("No profile selected")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	profile := f.GetProfile(*profileName)
	if profile == nil {
		err := errors.New("No profile selected")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}
	lockfile, err := f.GetCurrentLockfile(selectedInstall)
	if err != nil {
		err = errors.Wrap(err, "Failed to get lockfile")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}
	metadata := &ExportedProfileMetadata{
		GameVersion: selectedInstall.Info.Version,
	}

	exportedProfile := &ExportedProfile{
		Profile:  profile,
		LockFile: lockfile,
		Metadata: metadata,
	}

	defaultFileName := fmt.Sprintf("%s-%s.smmprofile", profile.Name, time.Now().UTC().Format("2006-01-02-15-04-05"))
	filename, err := wailsRuntime.SaveFileDialog(f.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: defaultFileName,
		Filters: []wailsRuntime.FileFilter{
			{
				Pattern:     "*.smmprofile",
				DisplayName: "SMM Profile (*.smmprofile)",
			},
		},
	})
	if err != nil {
		err = errors.Wrap(err, "Failed to open save dialog")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	exportedProfileJson, err := json.MarshalIndent(exportedProfile, "", "  ")
	if err != nil {
		err = errors.Wrap(err, "failed to marshal exported profile")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}
	err = os.WriteFile(filename, exportedProfileJson, 0755)
	if err != nil {
		err = errors.Wrap(err, "failed to write exported profile")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	if err != nil {
		err := errors.Wrapf(err, "Failed to save exported profile: %s", *profileName)
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	return nil
}

func (f *FicsitCLI) ReadExportedProfileMetadata(file string) (*ExportedProfileMetadata, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		err = errors.Wrap(err, "Failed to read exported profile")
		wailsRuntime.LogError(f.ctx, err.Error())
		return nil, err
	}

	var exportedProfile ExportedProfile
	err = json.Unmarshal(fileBytes, &exportedProfile)
	if err != nil {
		err = errors.Wrap(err, "Failed to unmarshal exported profile")
		wailsRuntime.LogError(f.ctx, err.Error())
		return nil, err
	}

	return exportedProfile.Metadata, nil
}

func (f *FicsitCLI) ImportProfile(name string, file string) error {
	selectedInstall := f.GetSelectedInstall()
	if selectedInstall == nil {
		err := errors.New("No installation selected")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	profileData, err := os.ReadFile(file)
	if err != nil {
		err = errors.Wrap(err, "Failed to read profile file")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	var exportedProfile ExportedProfile
	err = json.Unmarshal(profileData, &exportedProfile)
	if err != nil {
		err = errors.Wrap(err, "Failed to unmarshal profile file")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	profile, err := f.ficsitCli.Profiles.AddProfile(name)
	if err != nil {
		err = errors.Wrap(err, "Failed to add profile")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	profile.Mods = exportedProfile.Profile.Mods

	selectedInstall.Installation.SetProfile(f.ficsitCli, name)

	err = selectedInstall.Installation.WriteLockFile(f.ficsitCli, *exportedProfile.LockFile)
	if err != nil {
		f.ficsitCli.Profiles.DeleteProfile(name)
		err = errors.Wrap(err, "Failed to write lockfile")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()

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
		err = errors.Wrap(err, "Failed to get current lockfile")
		wailsRuntime.LogError(f.ctx, err.Error())
		return nil, err
	}

	if currentLockfile == nil {
		return nil, nil
	}

	profile := f.GetProfile(f.selectedInstallation.Installation.Profile)

	resolver := cli.NewDependencyResolver(f.ficsitCli.APIClient)

	gameVersion, err := f.selectedInstallation.Installation.GetGameVersion(f.ficsitCli)
	if err != nil {
		err = errors.Wrap(err, "failed to detect game version")
		wailsRuntime.LogError(f.ctx, err.Error())
		return nil, err
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
		err = errors.Wrap(err, "error resolving dependencies")
		wailsRuntime.LogError(f.ctx, err.Error())
		return nil, err
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
		err = errors.Wrap(err, "Failed to get current lockfile")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
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
		err = errors.Wrap(err, "Failed to update mods")
		wailsRuntime.LogError(f.ctx, err.Error())
		return err
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
	wailsRuntime.EventsEmit(f.ctx, "selectedProfile", profileName)
}

func (f *FicsitCLI) SetProgress(progress *Progress) {
	f.progress = progress
	wailsRuntime.EventsEmit(f.ctx, "progress", f.progress)
}
