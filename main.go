package main

import (
	"context"
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
	"github.com/satisfactorymodding/SatisfactoryModManager/autoupdate"
	"github.com/satisfactorymodding/SatisfactoryModManager/bindings"
	"github.com/satisfactorymodding/SatisfactoryModManager/file_scheme_association"
	"github.com/satisfactorymodding/SatisfactoryModManager/project_file"
	"github.com/satisfactorymodding/SatisfactoryModManager/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/singleinstance"
	"github.com/satisfactorymodding/SatisfactoryModManager/utils"
	"github.com/satisfactorymodding/SatisfactoryModManager/wails_logging"
	"github.com/satisfactorymodding/SatisfactoryModManager/websocket"
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
	if !singleinstance.RequestSingleInstanceLock() {
		return
	}

	autoupdate.Init(autoupdate.AutoUpdateConfig{
		UpdateFoundCallback: func(latestVersion string, changelogs map[string]string) {
			bindings.BindingsInstance.Update.UpdateAvailable(latestVersion, changelogs)
		},
		DownloadProgressCallback: func(bytesDownloaded, totalBytes int64) {
			bindings.BindingsInstance.Update.UpdateDownloadProgress(bytesDownloaded, totalBytes)
		},
		UpdateReadyCallback: func() { bindings.BindingsInstance.Update.UpdateReady() },
	})

	singleinstance.OnSecondInstance = func(args []string) {
		processArguments(args)
	}
	go singleinstance.ListenForSecondInstance()

	err := file_scheme_association.SetAsDefaultSchemeHandler("smmanager")
	if err != nil {
		log.Error().Err(err).Msg("Failed to set as default scheme handler")
	}

	err = file_scheme_association.SetAsDefaultFileHandler(".smmprofile")
	if err != nil {
		log.Error().Err(err).Msg("Failed to set as default file extension handler")
	}

	go websocket.ListenAndServeWebsocket()

	err = settings.LoadSettings()
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
		Width:     settings.Settings.UnexpandedSize.Width,
		Height:    settings.Settings.UnexpandedSize.Height,
		MinWidth:  utils.UnexpandedMin.Width,
		MaxWidth:  utils.UnexpandedMax.Width,
		MinHeight: utils.UnexpandedMin.Height,
		MaxHeight: utils.UnexpandedMax.Height,
		Assets:    assets,
		OnStartup: b.Startup,
		OnDomReady: func(ctx context.Context) {
			processArguments(os.Args)
			autoupdate.CheckInterval(5 * time.Minute)
		},
		Bind:   b.GetBindings(),
		Logger: wails_logging.WailsZeroLogLogger{},
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to start application")
	}

	err = autoupdate.OnExit(bindings.BindingsInstance.Update.Restart)
	if err != nil {
		log.Error().Err(err).Msg("Failed to apply update on exit")
	}
}

func init() {
	err := loadProjectFile()
	if err != nil {
		panic(err)
	}

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

	viper.Set("websocket-port", 33642)

	viper.Set("version", project_file.ProjectFile.Info.ProductVersion)
	viper.Set("github-release-repo", "satisfactorymodding/SatisfactoryModManager")

	// ficsit-cli config
	viper.Set("profiles-file", "profiles.json")
	viper.Set("installations-file", "installations.json")
	viper.Set("api-base", "https://api.ficsit.app")
	viper.Set("graphql-api", "/v2/query")

	viper.Set("log", "debug")
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
