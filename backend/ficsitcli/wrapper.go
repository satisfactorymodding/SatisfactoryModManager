package ficsitcli

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/satisfactorymodding/ficsit-cli/cli/provider"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	appCommon "github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
)

type ficsitCLI struct {
	ficsitCli            *cli.GlobalContext
	installationMetadata map[string]*common.Installation
	installFindErrors    []error
	progress             *Progress
	isGameRunning        bool
}

var FicsitCLI = &ficsitCLI{}

func (f *ficsitCLI) Init() error {
	if f.ficsitCli != nil {
		return fmt.Errorf("FicsitCLIWrapper already initialized")
	}
	var err error
	f.ficsitCli, err = cli.InitCLI(false)
	if err != nil {
		return fmt.Errorf("failed to initialize ficsit-cli: %w", err)
	}
	f.ficsitCli.Provider.(*provider.MixedProvider).Offline = settings.Settings.Offline
	err = f.initInstallations()
	if err != nil {
		return fmt.Errorf("failed to initialize installations: %w", err)
	}
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
	return nil
}

func (f *ficsitCLI) setProgress(p *Progress) {
	f.progress = p
	wailsRuntime.EventsEmit(appCommon.AppContext, "progress", p)
}

func (f *ficsitCLI) isValidInstall(path string) bool {
	_, ok := f.installationMetadata[path]
	return ok
}
