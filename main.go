package main

import (
	"embed"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/bindings"
	"github.com/satisfactorymodding/SatisfactoryModManager/utils"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	b, err := bindings.MakeBindings()
	if err != nil {
		panic(err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "SatisfactoryModManager",
		Frameless: runtime.GOOS == "windows",
		Width:     utils.UnexpandedMinWidth,
		Height:    b.Settings.Data.AppHeight,
		MinWidth:  utils.UnexpandedMinWidth,
		MaxWidth:  utils.UnexpandedMinWidth,
		MinHeight: utils.UnexpandedMinHeight,
		Assets:    assets,
		OnStartup: b.Startup,
		Bind:      b.GetBindings(),
	})

	if err != nil {
		println("Error:", err)
	}
}

func init() {
	// general config

	var baseLocalDir string

	switch runtime.GOOS {
	case "windows":
		baseLocalDir = os.Getenv("APPDATA")
	case "linux":
		baseLocalDir = path.Join(os.Getenv("HOME"), ".local", "share")
	default:
		panic("unsupported platform: " + runtime.GOOS)
	}

	viper.Set("base-local-dir", baseLocalDir)

	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	viper.Set("log", "info")

	cacheDir := filepath.Clean(filepath.Join(baseCacheDir, "SatisfactoryModManagerNEW"))
	utils.EnsureDirExists(cacheDir)
	viper.Set("cache-dir", cacheDir)

	localDir := filepath.Clean(filepath.Join(baseLocalDir, "SatisfactoryModManagerNEW"))
	utils.EnsureDirExists(localDir)
	viper.Set("local-dir", localDir)

	// ficsit-cli config
	viper.Set("profiles-file", "profiles.json")
	viper.Set("installations-file", "installations.json")
	viper.Set("api-base", "https://api.ficsit.app")
	viper.Set("graphql-api", "/v2/query")
}
