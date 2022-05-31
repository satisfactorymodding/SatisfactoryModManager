package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ExpandMod() bool {
	_, height := runtime.WindowGetSize(a.ctx)
	runtime.WindowSetMinSize(a.ctx, expandedMinWidth, expandedMinHeight)
	runtime.WindowSetMaxSize(a.ctx, -1, -1)
	runtime.WindowSetSize(a.ctx, expandedMinWidth, height)
	return true
}

func (a *App) UnexpandMod() bool {
	_, height := runtime.WindowGetSize(a.ctx)
	runtime.WindowSetMinSize(a.ctx, unexpandedMinWidth, unexpandedMinHeight)
	runtime.WindowSetMaxSize(a.ctx, unexpandedMinWidth, -1)
	runtime.WindowSetSize(a.ctx, unexpandedMinWidth, height)
	return true
}
