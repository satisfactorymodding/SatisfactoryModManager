package main

import (
	"embed"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed frontend/build
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

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
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err)
	}
}
