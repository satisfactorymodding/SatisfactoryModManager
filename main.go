package main

import (
	_ "embed"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/SatisfactoryModManager/bindings"
	"github.com/satisfactorymodding/SatisfactoryModManager/project_file"
	"github.com/satisfactorymodding/SatisfactoryModManager/utils"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"gopkg.in/natefinch/lumberjack.v2"
)

//go:embed wails.json
var projectFile []byte

func loadProjectFile() error {
	return json.Unmarshal(projectFile, &project_file.ProjectFile)
}

func main() {
	err := loadProjectFile()
	if err != nil {
		panic(err)
	}

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

	viper.Set("log", "info")
	viper.Set("log-file", filepath.Join(cacheDir, "logs", "SatisfactoryModManager.log"))

	writers := make([]io.Writer, 0)
	writers = append(writers, zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	if viper.GetString("log-file") != "" {
		logFile := &lumberjack.Logger{
			Filename:   viper.GetString("log-file"),
			MaxSize:    10, // megabytes
			MaxBackups: 5,
			MaxAge:     30, // days
		}

		if err == nil {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        logFile,
				TimeFormat: time.RFC3339,
				NoColor:    true,
			})
		}
	}

	log.Logger = zerolog.New(io.MultiWriter(writers...)).With().Timestamp().Logger()
}
