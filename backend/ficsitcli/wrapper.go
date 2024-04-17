package ficsitcli

import (
	"fmt"
	"log/slog"
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
	progress             *Progress
	isGameRunning        bool
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
	return nil
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
				if process.Executable() == "FactoryGame-Win64-Shipping.exe" || process.Executable() == "FactoryGame-Win64-Shipping" {
					f.isGameRunning = true
					break
				}
			}
			wailsRuntime.EventsEmit(appCommon.AppContext, "isGameRunning", f.isGameRunning)
		}
	}()
}

func (f *ficsitCLI) setProgress(p *Progress) {
	f.progress = p
	wailsRuntime.EventsEmit(appCommon.AppContext, "progress", p)
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
