package main

import (
	"context"
	"embed"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := &App{}
	ficsitCli := &FicsitCLI{}

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "SatisfactoryModManager",
		Frameless: runtime.GOOS == "windows",
		Width:     unexpandedMinWidth,
		Height:    unexpandedMinHeight,
		MaxWidth:  unexpandedMinWidth,
		MinWidth:  unexpandedMinWidth,
		MinHeight: unexpandedMinHeight,
		Assets:    assets,
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			ficsitCli.startup(ctx)
		},
		Bind: []interface{}{
			app,
			ficsitCli,
		},
	})

	if err != nil {
		println("Error:", err)
	}
}
