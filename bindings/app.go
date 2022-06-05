package bindings

import (
	"context"

	"github.com/satisfactorymodding/SatisfactoryModManager/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
}

func (a *App) ExpandMod() bool {
	_, height := runtime.WindowGetSize(a.ctx)
	runtime.WindowSetMinSize(a.ctx, utils.ExpandedMinWidth, utils.ExpandedMinHeight)
	runtime.WindowSetMaxSize(a.ctx, -1, -1)
	runtime.WindowSetSize(a.ctx, BindingsInstance.Settings.data.ExpandedAppWidth, height)
	a.isExpanded = true
	return true
}

func (a *App) UnexpandMod() bool {
	_, height := runtime.WindowGetSize(a.ctx)
	runtime.WindowSetMinSize(a.ctx, utils.UnexpandedMinWidth, utils.UnexpandedMinHeight)
	runtime.WindowSetMaxSize(a.ctx, utils.UnexpandedMinWidth, -1)
	runtime.WindowSetSize(a.ctx, utils.UnexpandedMinWidth, height)
	a.isExpanded = false
	return true
}
