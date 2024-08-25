package ficsitcli

import (
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/satisfactorymodding/ficsit-cli/cli/provider"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	appCommon "github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
)

type ficsitCLI struct {
	ficsitCli            *cli.GlobalContext
	installationMetadata *xsync.MapOf[string, installationMetadata]
	installFindErrors    []error
	isGameRunning        bool
	actionMutex          sync.Mutex
}

var FicsitCLI *ficsitCLI

func Init() error {
	if FicsitCLI != nil {
		return nil
	}
	ficsitCli, err := cli.InitCLI(false)
	if err != nil {
		return fmt.Errorf("failed to initialize ficsit-cli: %w", err)
	}
	ficsitCli.Provider.(*provider.MixedProvider).Offline = settings.Settings.Offline

	FicsitCLI = &ficsitCLI{ficsitCli: ficsitCli, installationMetadata: xsync.NewMapOf[string, installationMetadata]()}
	err = FicsitCLI.initInstallations()
	if err != nil {
		return fmt.Errorf("failed to initialize installations: %w", err)
	}

	if settings.SMM2SelectedProfile != nil {
		for _, install := range FicsitCLI.ficsitCli.Installations.Installations {
			profile := settings.SMM2SelectedProfile[install.Path]
			if profile != "" {
				err := install.SetProfile(FicsitCLI.ficsitCli, profile)
				if err != nil {
					slog.Error(
						"failed to restore selected profile, using fallback",
						slog.String("install", install.Path),
						slog.String("profile", profile),
						slog.Any("error", err),
					)
					install.Profile = FicsitCLI.GetFallbackProfile()
				}
			}
		}
	}

	return nil
}

// With and without `.exe` variants in case it is missing on Linux
var executableNames = []string{
	"FactoryGame-Win64-Shipping.exe", "FactoryGame-Win64-Shipping",
	"FactoryGameSteam-Win64-Shipping.exe", "FactoryGameSteam-Win64-Shipping",
	"FactoryGameEGS-Win64-Shipping.exe", "FactoryGameEGS-Win64-Shipping",
}

func (f *ficsitCLI) StartGameRunningWatcher() {
	gameRunningTicker := time.NewTicker(5 * time.Second)
	go func() {
		for range gameRunningTicker.C {
			processes, err := ps.Processes()
			if err != nil {
				slog.Error("failed to get processes", slog.Any("error", err))
				continue
			}
			f.isGameRunning = false
			for _, process := range processes {
				if slices.Contains(executableNames, process.Executable()) {
					f.isGameRunning = true
					break
				}
			}
			wailsRuntime.EventsEmit(appCommon.AppContext, "isGameRunning", f.isGameRunning)
		}
	}()
}

// GetProgress exists only to ensure the Progress type is exported to typescript. It returns nil
func (f *ficsitCLI) GetProgress() *Progress {
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

	installsWithTargets, _, err := f.getInstallsToApply()
	if err != nil {
		slog.Error("failed to get installs to apply", slog.Any("error", err))
		return
	}
	installsForTarget := make(map[string][]string)
	for _, install := range installsWithTargets {
		installsForTarget[install.targetName] = append(installsForTarget[install.targetName], install.install.Path)
	}
	wailsRuntime.EventsEmit(appCommon.AppContext, "selectedProfileTargets", installsForTarget)
}

func (f *ficsitCLI) isValidInstall(path string) bool {
	meta, ok := f.installationMetadata.Load(path)
	return ok && meta.State != InstallStateInvalid
}

func (f *ficsitCLI) WipeMods(includeRemote bool) error {
	for _, i := range f.ficsitCli.Installations.Installations {
		if !includeRemote {
			meta, ok := f.installationMetadata.Load(i.Path)
			if !ok {
				// If the metadata is not registered yet, it is definitely not a local installation
				continue
			}
			if meta.Info == nil {
				// If the metadata is not available, it is definitely not a local installation
				continue
			}
			if meta.Info.Location != common.LocationTypeLocal {
				continue
			}
		}

		err := i.Wipe()
		if err != nil {
			return fmt.Errorf("failed to wipe installation %s: %w", i.Path, err)
		}
	}
	return nil
}
