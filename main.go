package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/satisfactorymodding/SatisfactoryModManager/association"
	"github.com/satisfactorymodding/SatisfactoryModManager/autoupdate"
	"github.com/satisfactorymodding/SatisfactoryModManager/bindings"
	"github.com/satisfactorymodding/SatisfactoryModManager/projectfile"
	"github.com/satisfactorymodding/SatisfactoryModManager/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/singleinstance"
	"github.com/satisfactorymodding/SatisfactoryModManager/utils"
	"github.com/satisfactorymodding/SatisfactoryModManager/wailslogging"
	"github.com/satisfactorymodding/SatisfactoryModManager/websocket"

	_ "embed"
)

//go:embed wails.json
var projectFile []byte

func loadProjectFile() error {
	err := json.Unmarshal(projectFile, &projectfile.ProjectFile)
	if err != nil {
		return errors.Wrap(err, "Failed to load project file")
	}
	return nil
}

func main() {
	if !singleinstance.RequestSingleInstanceLock() {
		return
	}

	autoupdate.Init(autoupdate.Config{
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

	err := association.SetAsDefaultSchemeHandler("smmanager")
	if err != nil {
		log.Error().Err(err).Msg("Failed to set as default scheme handler")
	}

	err = association.SetAsDefaultFileHandler(".smmprofile")
	if err != nil {
		log.Error().Err(err).Msg("Failed to set as default file extension handler")
	}

	go websocket.ListenAndServeWebsocket()

	err = settings.LoadSettings()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load settings")
	}

	b, err := bindings.MakeBindings()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create bindings")
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
		Logger: wailslogging.WailsZeroLogLogger{},
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
		log.Fatal().Str("os", runtime.GOOS).Msg("Unsupported OS")
	}

	viper.Set("base-local-dir", baseLocalDir)

	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	cacheDir := filepath.Clean(filepath.Join(baseCacheDir, "SatisfactoryModManagerNEW"))
	_ = utils.EnsureDirExists(cacheDir)
	viper.Set("cache-dir", cacheDir)

	localDir := filepath.Clean(filepath.Join(baseLocalDir, "SatisfactoryModManagerNEW"))
	_ = utils.EnsureDirExists(localDir)
	viper.Set("local-dir", localDir)

	viper.Set("websocket-port", 33642)

	viper.Set("version", projectfile.ProjectFile.Info.ProductVersion)
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

		writers = append(writers, zerolog.ConsoleWriter{
			Out:        logFile,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		})
	}

	level, err := zerolog.ParseLevel(viper.GetString("log"))
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = zerolog.New(io.MultiWriter(writers...)).With().Timestamp().Logger()
}
