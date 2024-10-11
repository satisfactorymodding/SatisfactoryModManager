package ficsitcli

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	resolver "github.com/satisfactorymodding/ficsit-resolver"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

type taskUpdate struct {
	taskName string
	progress utils.Progress
}

func (f *ficsitCLI) action(action Action, item ProgressItem, run func(*slog.Logger, chan<- taskUpdate) error) error {
	if !f.actionMutex.TryLock() {
		return fmt.Errorf("another operation in progress")
	}
	defer f.actionMutex.Unlock()

	var logAttrs []any
	logAttrs = append(logAttrs, slog.String("type", string(action)))
	if item != noItem {
		logAttrs = append(logAttrs, slog.String("item", item.Name))
		if item.Version != "" {
			logAttrs = append(logAttrs, slog.String("version", item.Version))
		}
	}
	l := slog.With(slog.Group("action", logAttrs...))

	done := make(chan bool)
	defer close(done)

	progress := newProgress(action, item)
	tasks := xsync.NewMapOf[string, utils.Progress]()
	go func() {
		wailsRuntime.EventsEmit(common.AppContext, "progress", progress)
		defer wailsRuntime.EventsEmit(common.AppContext, "progress", nil)

		progressTicker := time.NewTicker(100 * time.Millisecond)
		defer progressTicker.Stop()

		for {
			select {
			case <-done:
				return
			case <-progressTicker.C:
				tasks.Range(func(key string, value utils.Progress) bool {
					progress.Tasks[key] = value
					return true
				})
				wailsRuntime.EventsEmit(common.AppContext, "progress", progress)
			}
		}
	}()

	taskChannel := make(chan taskUpdate)
	go func() {
		for update := range taskChannel {
			tasks.LoadAndStore(update.taskName, update.progress)
		}
	}()

	err := run(l, taskChannel)
	if err != nil {
		l.Info("action failed")
		return err
	}

	l.Info("action complete")

	return nil
}

func (f *ficsitCLI) Apply() error {
	profileName := f.GetSelectedProfile()
	if profileName == nil {
		return fmt.Errorf("no profile selected")
	}
	return f.action(ActionApply, newSimpleItem(*profileName), f.apply)
}

func (f *ficsitCLI) apply(l *slog.Logger, taskChannel chan<- taskUpdate) error {
	installsToApply, profile, err := f.getInstallsToApply()
	if err != nil {
		return err
	}

	targetsUsingProfile := make(map[resolver.TargetName]bool)
	for _, install := range installsToApply {
		targetsUsingProfile[resolver.TargetName(install.targetName)] = true
	}

	profile.RequiredTargets = maps.Keys(targetsUsingProfile)
	err = f.ficsitCli.Profiles.Save()
	if err != nil {
		l.Error("failed to save profile", slog.Any("error", err))
	}

	f.EmitModsChange()
	defer f.EmitModsChange()

	defer close(taskChannel)

	var errg errgroup.Group
	var wg sync.WaitGroup

	for _, installTarget := range installsToApply {
		wg.Add(1)
		errg.Go(func() error {
			defer wg.Done()

			installChannel := make(chan cli.InstallUpdate)

			go func() {
				for update := range installChannel {
					switch update.Type {
					case cli.InstallUpdateTypeModDownload:
						taskChannel <- taskUpdate{
							taskName: fmt.Sprintf("%s:%s:%s:download", update.Item.Mod, update.Item.Version, installTarget.targetName),
							progress: utils.Progress{
								Current: update.Progress.Completed,
								Total:   update.Progress.Total,
							},
						}
					case cli.InstallUpdateTypeModExtract:
						taskChannel <- taskUpdate{
							taskName: fmt.Sprintf("%s:%s:%s:extract", update.Item.Mod, update.Item.Version, installTarget.targetName),
							progress: utils.Progress{
								Current: update.Progress.Completed,
								Total:   update.Progress.Total,
							},
						}
					}
				}
			}()

			installErr := installTarget.install.Install(f.ficsitCli, installChannel)
			if installErr != nil {
				var solvingError resolver.DependencyResolverError
				if errors.As(installErr, &solvingError) {
					return solvingError
				}
				return installErr //nolint:wrapcheck
			}
			return nil
		})
	}

	if err := errg.Wait(); err != nil {
		// Ensure everything is finished, but return first error
		wg.Wait()
		return err //nolint:wrapcheck
	}

	return nil
}

type installWithTarget struct {
	install    *cli.Installation
	targetName string
}

func (f *ficsitCLI) getInstallsToApply() ([]installWithTarget, *cli.Profile, error) {
	selectedInstall := f.GetSelectedInstall()
	if selectedInstall == nil {
		return nil, nil, fmt.Errorf("no installation selected")
	}

	selectedProfile := selectedInstall.Profile
	allInstalls := f.GetInstallations()

	var installsUsingProfile []installWithTarget
	targetsUsingProfile := make(map[resolver.TargetName]bool)
	for _, install := range allInstalls {
		meta, ok := f.installationMetadata.Load(install)
		if !ok {
			continue
		}
		if meta.State != InstallStateValid {
			continue
		}
		i := f.GetInstallation(install)
		if i.Profile == selectedProfile {
			platform, err := i.GetPlatform(f.ficsitCli)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get platform: %w", err)
			}
			installsUsingProfile = append(installsUsingProfile, installWithTarget{
				install:    i,
				targetName: platform.TargetName,
			})
			targetsUsingProfile[resolver.TargetName(platform.TargetName)] = true
		}
	}

	return installsUsingProfile, f.GetProfile(selectedProfile), nil
}
