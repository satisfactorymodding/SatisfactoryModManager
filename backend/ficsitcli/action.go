package ficsitcli

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	resolver "github.com/satisfactorymodding/ficsit-resolver"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
)

type taskUpdate struct {
	taskName string
	progress ProgressTask
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
	tasks := xsync.NewMapOf[string, ProgressTask]()
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
				tasks.Range(func(key string, value ProgressTask) bool {
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

func (f *ficsitCLI) validateInstall(installation *cli.Installation, taskChannel chan<- taskUpdate) error {
	if !f.isValidInstall(installation.Path) {
		return fmt.Errorf("invalid installation: %s", installation.Path)
	}

	f.EmitModsChange()
	defer f.EmitModsChange()

	installChannel := make(chan cli.InstallUpdate)

	go func() {
		defer close(taskChannel)
		for update := range installChannel {
			switch update.Type {
			case cli.InstallUpdateTypeModDownload:
				taskChannel <- taskUpdate{
					taskName: fmt.Sprintf("%s:download", update.Item.Mod),
					progress: ProgressTask{
						Current: update.Progress.Completed,
						Total:   update.Progress.Total,
					},
				}
			case cli.InstallUpdateTypeModExtract:
				taskChannel <- taskUpdate{
					taskName: fmt.Sprintf("%s:extract", update.Item.Mod),
					progress: ProgressTask{
						Current: update.Progress.Completed,
						Total:   update.Progress.Total,
					},
				}
			}
		}
	}()

	installErr := installation.Install(f.ficsitCli, installChannel)
	if installErr != nil {
		var solvingError resolver.DependencyResolverError
		if errors.As(installErr, &solvingError) {
			return solvingError
		}
		return installErr //nolint:wrapcheck
	}
	return nil
}
