package bindings

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

// App struct
type App struct {
	ctx             context.Context
	isExpanded      bool
	stopSizeWatcher chan bool
}

func MakeApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.stopSizeWatcher = make(chan bool)
	windowStateTicker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-a.stopSizeWatcher:
				return
			case <-windowStateTicker.C:
				if wailsRuntime.WindowIsMinimised(a.ctx) {
					// When the window is minimized, the window position values are garbage
					continue
				}
				w, h := wailsRuntime.WindowGetSize(a.ctx)
				x, y := wailsRuntime.WindowGetPosition(a.ctx)
				changed := false
				if BindingsInstance.App.isExpanded {
					if w != settings.Settings.ExpandedSize.Width {
						settings.Settings.ExpandedSize.Width = w
						changed = true
					}
				} else {
					if w != settings.Settings.UnexpandedSize.Width {
						settings.Settings.UnexpandedSize.Width = w
						changed = true
					}
				}
				if h != settings.Settings.ExpandedSize.Height {
					settings.Settings.ExpandedSize.Height = h
					changed = true
				}
				if h != settings.Settings.UnexpandedSize.Height {
					settings.Settings.UnexpandedSize.Height = h
					changed = true
				}
				if settings.Settings.WindowPosition == nil {
					settings.Settings.WindowPosition = &utils.Position{
						X: x,
						Y: y,
					}
					changed = true
				} else {
					if x != settings.Settings.WindowPosition.X {
						settings.Settings.WindowPosition.X = x
						changed = true
					}
					if y != settings.Settings.WindowPosition.Y {
						settings.Settings.WindowPosition.Y = y
						changed = true
					}
				}
				if changed {
					err := settings.SaveSettings()
					if err != nil {
						slog.Error("failed to save settings", slog.Any("error", err))
					}
				}
			}
		}
	}()

	// Wails doesn't support setting the window position on init, so we do it here
	if settings.Settings.WindowPosition != nil {
		wailsRuntime.WindowSetPosition(ctx, settings.Settings.WindowPosition.X, settings.Settings.WindowPosition.Y)
	}
}

func (a *App) ExpandMod() bool {
	_, height := wailsRuntime.WindowGetSize(a.ctx)
	wailsRuntime.WindowSetMinSize(a.ctx, utils.ExpandedMin.Width, utils.ExpandedMin.Height)
	wailsRuntime.WindowSetMaxSize(a.ctx, utils.ExpandedMax.Width, utils.ExpandedMax.Height)
	wailsRuntime.WindowSetSize(a.ctx, settings.Settings.ExpandedSize.Width, height)
	a.isExpanded = true
	return true
}

func (a *App) UnexpandMod() bool {
	a.isExpanded = false
	_, height := wailsRuntime.WindowGetSize(a.ctx)
	wailsRuntime.WindowSetMinSize(a.ctx, utils.UnexpandedMin.Width, utils.UnexpandedMin.Height)
	wailsRuntime.WindowSetMaxSize(a.ctx, utils.UnexpandedMax.Width, utils.UnexpandedMax.Height)
	wailsRuntime.WindowSetSize(a.ctx, settings.Settings.UnexpandedSize.Width, height)
	return true
}

func (a *App) GetVersion() string {
	return viper.GetString("version")
}

func (a *App) GetCommit() string {
	return viper.GetString("commit")
}

func (a *App) GetDate() string {
	return viper.GetString("date")
}

type FileFilter struct {
	DisplayName string `json:"displayName"`
	Pattern     string `json:"pattern"`
}

type OpenDialogOptions struct {
	DefaultDirectory           string       `json:"defaultDirectory,omitempty"`
	DefaultFilename            string       `json:"defaultFilename,omitempty"`
	Title                      string       `json:"title,omitempty"`
	Filters                    []FileFilter `json:"filters,omitempty"`
	ShowHiddenFiles            bool         `json:"showHiddenFiles,omitempty"`
	CanCreateDirectories       bool         `json:"canCreateDirectories,omitempty"`
	ResolvesAliases            bool         `json:"resolvesAliases,omitempty"`
	TreatPackagesAsDirectories bool         `json:"treatPackagesAsDirectories,omitempty"`
}

func (a *App) OpenFileDialog(options OpenDialogOptions) (string, error) {
	wailsFilters := make([]wailsRuntime.FileFilter, len(options.Filters))
	for i, filter := range options.Filters {
		wailsFilters[i] = wailsRuntime.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		}
	}
	wailsOptions := wailsRuntime.OpenDialogOptions{
		DefaultDirectory:           options.DefaultDirectory,
		DefaultFilename:            options.DefaultFilename,
		Title:                      options.Title,
		Filters:                    wailsFilters,
		ShowHiddenFiles:            options.ShowHiddenFiles,
		CanCreateDirectories:       options.CanCreateDirectories,
		ResolvesAliases:            options.ResolvesAliases,
		TreatPackagesAsDirectories: options.TreatPackagesAsDirectories,
	}
	file, err := wailsRuntime.OpenFileDialog(a.ctx, wailsOptions)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file dialog")
	}
	return file, nil
}

func (a *App) OpenDirectoryDialog(options OpenDialogOptions) (string, error) {
	wailsFilters := make([]wailsRuntime.FileFilter, len(options.Filters))
	for i, filter := range options.Filters {
		wailsFilters[i] = wailsRuntime.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		}
	}
	wailsOptions := wailsRuntime.OpenDialogOptions{
		DefaultDirectory:           options.DefaultDirectory,
		DefaultFilename:            options.DefaultFilename,
		Title:                      options.Title,
		Filters:                    wailsFilters,
		ShowHiddenFiles:            options.ShowHiddenFiles,
		CanCreateDirectories:       options.CanCreateDirectories,
		ResolvesAliases:            options.ResolvesAliases,
		TreatPackagesAsDirectories: options.TreatPackagesAsDirectories,
	}
	file, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsOptions)
	if err != nil {
		return "", errors.Wrap(err, "failed to open directory dialog")
	}
	return file, nil
}

func (a *App) ExternalInstallMod(modID, version string) {
	wailsRuntime.EventsEmit(a.ctx, "externalInstallMod", modID, version)
}

func (a *App) ExternalImportProfile(path string) {
	wailsRuntime.EventsEmit(a.ctx, "externalImportProfile", path)
}

func (a *App) Show() {
	wailsRuntime.WindowUnminimise(a.ctx)
	wailsRuntime.Show(a.ctx)
}

func (a *App) OpenExternal(input string) {
	err := browser.OpenFile(input)
	if err != nil {
		slog.Error("failed to open external", slog.Any("error", err), slog.String("path", input))
	}
}

func (a *App) GetAPIEndpoint() string {
	return viper.GetString("api-base") + viper.GetString("graphql-api")
}

func (a *App) GetSiteEndpoint() string {
	return strings.Replace(viper.GetString("api-base"), "api.", "", 1)
}

func (a *App) shutdown(_ context.Context) {
	if a.stopSizeWatcher != nil {
		close(a.stopSizeWatcher)
	}
}
