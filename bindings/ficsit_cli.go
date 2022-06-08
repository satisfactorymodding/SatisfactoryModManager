package bindings

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/ficsit-cli/cli"

	"github.com/satisfactorymodding/SatisfactoryModManager/installFinders"
	installFinderTypes "github.com/satisfactorymodding/SatisfactoryModManager/installFinders/types"
)

type FicsitCLI struct {
	ctx                  context.Context
	ficsitCli            *cli.GlobalContext
	installations        []*InstallationInfo
	selectedInstallation *InstallationInfo
	progress             *Progress
}

type InstallationInfo struct {
	Installation *cli.Installation                `json:"installation"`
	Info         *installFinderTypes.Installation `json:"info"`
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

	f.initInstallations()

	return f, nil
}

func (f *FicsitCLI) startup(ctx context.Context) {
	f.ctx = ctx
}

func (f *FicsitCLI) initInstallations() {
	installs, _, findErrors := installFinders.FindInstallations()

	if len(findErrors) > 0 {
		for _, err := range findErrors {
			wailsRuntime.LogDebugf(f.ctx, "Failed to find installations: %v", err)
		}
	}

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

func (f *FicsitCLI) SetProfile(profile string) {
	f.GetInstallation(f.selectedInstallation.Info.Path).Installation.SetProfile(f.ficsitCli, profile)
	f.emitModsChange()
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

func (f *FicsitCLI) InstallMod(mod string) {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)

	profileErr := profile.AddMod(mod, ">=0.0.0")
	if profileErr != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to add mod: %v", profileErr)
		return
	}

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	})

	installChannel := make(chan cli.InstallUpdate)

	go func() {
		for data := range installChannel {
			if data.DownloadProgress < 1 {
				f.SetProgress(&Progress{
					Item:     mod,
					Message:  "Downloading " + data.ModName,
					Progress: data.DownloadProgress,
				})
			} else {
				f.SetProgress(&Progress{
					Item:     mod,
					Message:  "Extracting " + data.ModName,
					Progress: data.DownloadProgress,
				})
			}
		}
	}()

	installErr := installation.Installation.Install(f.ficsitCli, installChannel)

	close(installChannel)

	if installErr != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to install mod: %v", installErr)
		return
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()
	f.SetProgress(nil)
}

func (f *FicsitCLI) InstallModVersion(mod string, version string) {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)
	profileErr := profile.AddMod(mod, version)

	if profileErr != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to add mod: %v", profileErr)
		return
	}

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Finding the best version to install",
		Progress: -1,
	})

	installChannel := make(chan cli.InstallUpdate)

	go func() {
		for data := range installChannel {
			if data.DownloadProgress < 1 {
				f.SetProgress(&Progress{
					Item:     mod,
					Message:  "Downloading " + data.ModName,
					Progress: data.DownloadProgress,
				})
			} else {
				f.SetProgress(&Progress{
					Item:     mod,
					Message:  "Extracting " + data.ModName,
					Progress: data.DownloadProgress,
				})
			}
		}
	}()

	installErr := installation.Installation.Install(f.ficsitCli, installChannel)

	close(installChannel)

	if installErr != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to install mod: %v", installErr)
		return
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()
	f.SetProgress(nil)
}

func (f *FicsitCLI) RemoveMod(mod string) {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)
	profile.RemoveMod(mod)

	f.SetProgress(&Progress{
		Item:     mod,
		Message:  "Checking for mods that are no longer needed",
		Progress: -1,
	})

	installChannel := make(chan cli.InstallUpdate)

	go func() {
		for data := range installChannel {
			if data.DownloadProgress < 1 {
				f.SetProgress(&Progress{
					Item:     mod,
					Message:  "Downloading " + data.ModName,
					Progress: data.DownloadProgress,
				})
			} else {
				f.SetProgress(&Progress{
					Item:     mod,
					Message:  "Extracting " + data.ModName,
					Progress: data.ExtractProgress,
				})
			}
		}
	}()

	installErr := installation.Installation.Install(f.ficsitCli, installChannel)

	close(installChannel)

	if installErr != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to install mod: %v", installErr)
		return
	}

	f.ficsitCli.Profiles.Save()
	f.emitModsChange()
	f.SetProgress(nil)
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
