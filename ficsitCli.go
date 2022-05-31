package main

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/ficsit-cli/cli"

	"github.com/satisfactorymodding/SatisfactoryModManager/installFinders"
	installFinderTypes "github.com/satisfactorymodding/SatisfactoryModManager/installFinders/types"
)

func InitFicsitCLI() (*cli.GlobalContext, error) {
	var baseLocalDir string

	switch runtime.GOOS {
	case "windows":
		baseLocalDir = os.Getenv("APPDATA")
	case "linux":
		baseLocalDir = path.Join(os.Getenv("HOME"), ".local", "share")
	default:
		panic("unsupported platform: " + runtime.GOOS)
	}

	viper.Set("base-local-dir", baseLocalDir)

	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	viper.Set("log", "info")
	viper.Set("cache-dir", filepath.Clean(filepath.Join(baseCacheDir, "SatisfactoryModManagerNEW")))
	viper.Set("local-dir", filepath.Clean(filepath.Join(baseLocalDir, "SatisfactoryModManagerNEW")))
	viper.Set("profiles-file", "profiles.json")
	viper.Set("installations-file", "installations.json")
	viper.Set("api-base", "https://api.ficsit.app")
	viper.Set("graphql-api", "/v2/query")

	return cli.InitCLI()
}

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

func (f *FicsitCLI) startup(ctx context.Context) {
	f.ctx = ctx

	ficsitCli, err := InitFicsitCLI()
	if err != nil {
		wailsRuntime.LogErrorf(ctx, "Failed to initialize CLI: %v", err)
		return
	}
	f.ficsitCli = ficsitCli

	f.InitInstallations()
}

func (f *FicsitCLI) InitInstallations() {
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
	f.EmitModsChange()
}

func (f *FicsitCLI) SetProfile(profile string) {
	f.GetInstallation(f.selectedInstallation.Info.Path).Installation.SetProfile(f.ficsitCli, profile)
	f.EmitModsChange()
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

	installErr := installation.Installation.Install(f.ficsitCli)
	if installErr != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to install mod: %v", installErr)
		return
	}

	f.ficsitCli.Profiles.Save()
	f.EmitModsChange()
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

	installErr := installation.Installation.Install(f.ficsitCli)
	if installErr != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to install mod: %v", installErr)
		return
	}

	f.ficsitCli.Profiles.Save()
	f.EmitModsChange()
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

	installErr := installation.Installation.Install(f.ficsitCli)
	if installErr != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed to install mod: %v", installErr)
		return
	}

	f.ficsitCli.Profiles.Save()
	f.EmitModsChange()
	f.SetProgress(nil)
}

func (f *FicsitCLI) GetLockfile(i *cli.Installation) (*cli.LockFile, error) {
	platform, err := i.GetPlatform(f.ficsitCli)
	if err != nil {
		return nil, err
	}

	lockfilePath := path.Join(i.Path, platform.LockfilePath)

	var lockFile *cli.LockFile
	lockFileJSON, err := os.ReadFile(lockfilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.Wrap(err, "failed reading lockfile")
		}
	} else {
		if err := json.Unmarshal(lockFileJSON, &lockFile); err != nil {
			return nil, errors.Wrap(err, "failed parsing lockfile")
		}
	}

	return lockFile, nil
}

func (f *FicsitCLI) EmitModsChange() {
	installation := f.GetInstallation(f.selectedInstallation.Info.Path)
	profileName := installation.Installation.Profile
	profile := f.GetProfile(profileName)
	wailsRuntime.EventsEmit(f.ctx, "manifestMods", profile.Mods)
	lockfile, err := f.GetLockfile(installation.Installation)
	if err != nil {
		wailsRuntime.LogErrorf(f.ctx, "Failed read lockfile: %v", err)
		return
	}
	wailsRuntime.EventsEmit(f.ctx, "lockfileMods", lockfile)
}

func (f *FicsitCLI) SetProgress(progress *Progress) {
	f.progress = progress
	wailsRuntime.EventsEmit(f.ctx, "progress", f.progress)
}
