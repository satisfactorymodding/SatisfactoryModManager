package bindings

import (
	"context"
	"time"

	"github.com/satisfactorymodding/SatisfactoryModManager/project_file"
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
			if BindingsInstance.App.isExpanded {
				if w != BindingsInstance.Settings.Data.ExpandedAppWidth {
					BindingsInstance.Settings.Data.ExpandedAppWidth = w
					BindingsInstance.Settings.save()
				}
			}
			if h != BindingsInstance.Settings.Data.AppHeight {
				BindingsInstance.Settings.Data.AppHeight = h
				BindingsInstance.Settings.save()
			}
		}
	}()
}

func (a *App) ExpandMod() bool {
	_, height := wailsRuntime.WindowGetSize(a.ctx)
	wailsRuntime.WindowSetMinSize(a.ctx, utils.ExpandedMinWidth, utils.ExpandedMinHeight)
	wailsRuntime.WindowSetMaxSize(a.ctx, -1, -1)
	wailsRuntime.WindowSetSize(a.ctx, BindingsInstance.Settings.Data.ExpandedAppWidth, height)
	a.isExpanded = true
	return true
}

func (a *App) UnexpandMod() bool {
	_, height := wailsRuntime.WindowGetSize(a.ctx)
	wailsRuntime.WindowSetMinSize(a.ctx, utils.UnexpandedMinWidth, utils.UnexpandedMinHeight)
	wailsRuntime.WindowSetMaxSize(a.ctx, utils.UnexpandedMinWidth, -1)
	wailsRuntime.WindowSetSize(a.ctx, utils.UnexpandedMinWidth, height)
	a.isExpanded = false
	return true
}

func (a *App) GetVersion() string {
	return project_file.Version()
}
