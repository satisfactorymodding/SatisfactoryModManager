package bindings

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/SatisfactoryModManager/project_file"
	"github.com/satisfactorymodding/SatisfactoryModManager/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/utils"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
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

	sizeTicker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range sizeTicker.C {
			w, h := wailsRuntime.WindowGetSize(a.ctx)
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
			if changed {
				err := settings.SaveSettings()
				if err != nil {
					log.Error().Err(err).Msg("Failed to save settings")
				}
			}
		}
	}()
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
	return project_file.Version()
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
	return wailsRuntime.OpenFileDialog(a.ctx, wailsOptions)
}

func (a *App) ExternalInstallMod(modId, version string) {
	wailsRuntime.EventsEmit(a.ctx, "externalInstallMod", modId, version)
}

func (a *App) ExternalImportProfile(path string) {
	wailsRuntime.EventsEmit(a.ctx, "externalImportProfile", path)
}

func (a *App) Show() {
	wailsRuntime.Show(a.ctx)
}
