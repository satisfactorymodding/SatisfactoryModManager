package bindings

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/projectfile"
	"github.com/satisfactorymodding/SatisfactoryModManager/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/utils"
)

// App struct
type App struct {
	ctx        context.Context
	isExpanded bool
}

func MakeApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	windowStateTicker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range windowStateTicker.C {
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
					log.Error().Err(err).Msg("Failed to save settings")
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
	return projectfile.Version()
}

type FileFilter struct {
	DisplayName string `json:"displayName"`
	Pattern     string `json:"pattern"`
}

type OpenDialogOptions struct {
	DefaultDirectory           string       `json:"defaultDirectory"`
	DefaultFilename            string       `json:"defaultFilename"`
	Title                      string       `json:"title"`
	Filters                    []FileFilter `json:"filters"`
	ShowHiddenFiles            bool         `json:"showHiddenFiles"`
	CanCreateDirectories       bool         `json:"canCreateDirectories"`
	ResolvesAliases            bool         `json:"resolvesAliases"`
	TreatPackagesAsDirectories bool         `json:"treatPackagesAsDirectories"`
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

func (a *App) ExternalInstallMod(modID, version string) {
	wailsRuntime.EventsEmit(a.ctx, "externalInstallMod", modID, version)
}

func (a *App) ExternalImportProfile(path string) {
	wailsRuntime.EventsEmit(a.ctx, "externalImportProfile", path)
}

func (a *App) Show() {
	wailsRuntime.Show(a.ctx)
}

func (a *App) GetAPIEndpoint() string {
	return viper.GetString("api-base") + viper.GetString("graphql-api")
}

func (a *App) GetSiteEndpoint() string {
	return strings.Replace(viper.GetString("api-base"), "api.", "", 1)
}
