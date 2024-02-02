package ficsitcli

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/satisfactorymodding/ficsit-cli/cli/provider"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
)

type FicsitCLI struct {
	ctx                  context.Context
	ficsitCli            *cli.GlobalContext
	installationMetadata map[string]*common.Installation
	installFindErrors    []error
	progress             *Progress
	isGameRunning        bool
}

func MakeFicsitCLI() *FicsitCLI {
	return &FicsitCLI{}
}

func (f *FicsitCLI) Startup(ctx context.Context) error {
	f.ctx = ctx
	return f.init()
}

func (f *FicsitCLI) init() error {
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
			wailsRuntime.EventsEmit(f.ctx, "isGameRunning", f.isGameRunning)
		}
	}()
	return nil
}

func (f *FicsitCLI) setProgress(p *Progress) {
	f.progress = p
	wailsRuntime.EventsEmit(f.ctx, "progress", p)
}

func (f *FicsitCLI) isValidInstall(path string) bool {
	_, ok := f.installationMetadata[path]
	return ok
}
