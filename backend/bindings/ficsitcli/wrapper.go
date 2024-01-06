package ficsitcli

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/satisfactorymodding/ficsit-cli/cli/provider"
	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
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

func MakeFicsitCLI() (*FicsitCLI, error) {
	f := &FicsitCLI{}
	err := f.Init()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *FicsitCLI) Startup(ctx context.Context) {
	f.ctx = ctx
}

func (f *FicsitCLI) Init() error {
	if f.ficsitCli != nil {
		return errors.New("FicsitCLIWrapper already initialized")
	}
	var err error
	f.ficsitCli, err = cli.InitCLI(false)
	if err != nil {
		return errors.Wrap(err, "failed to initialize ficsit-cli")
	}
	f.ficsitCli.Provider.(*provider.MixedProvider).Offline = settings.Settings.Offline
	err = f.initInstallations()
	if err != nil {
		return errors.Wrap(err, "failed to initialize installations")
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

func ValidateCacheDir(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrapf(err, "failed to stat %s", dir)
		}
	} else {
		if !stat.IsDir() {
			return errors.Errorf("%s is not a directory", dir)
		}
	}
	return nil
}

func MoveCacheDir(newDir string) error {
	if newDir == viper.GetString("cache-dir") {
		return nil
	}

	err := ValidateCacheDir(newDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(newDir, 0o755)
	if err != nil {
		if !os.IsExist(err) {
			return errors.Wrapf(err, "failed to create %s", newDir)
		}
	}

	items, err := os.ReadDir(newDir)
	if err != nil {
		return errors.Wrapf(err, "failed to check if directory %s is empty", newDir)
	}
	if len(items) > 0 {
		return errors.Errorf("directory %s is not empty", newDir)
	}

	oldCacheDir := viper.GetString("cache-dir")
	// Move contents of oldCacheDir to dir
	if oldCacheDir != "" && oldCacheDir != newDir {
		err := moveCacheData(oldCacheDir, newDir)
		if err != nil {
			return err
		}
	}

	viper.Set("cache-dir", newDir)
	return nil
}

func moveCacheData(oldCacheDir, newDir string) error {
	oldStat, err := os.Stat(oldCacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Nothing to move
			return nil
		}
		return errors.Wrapf(err, "failed to stat %s", oldCacheDir)
	}
	if !oldStat.IsDir() {
		return errors.Errorf("%s is not a directory", oldCacheDir)
	}

	// Perform the move atomically
	copySuccess, err := utils.MoveRecursive(oldCacheDir, newDir)
	if err != nil {
		if !copySuccess {
			return err
		}
		slog.Error("failed to move cache dir", slog.Any("error", err))
	}

	return nil
}
