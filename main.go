package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/viper"
	"github.com/tawesoft/golib/v2/dialog"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/app"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/bindings"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/websocket"
)

//go:embed all:frontend/build
var assets embed.FS

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"

	updateMode = "none"
)

func main() {
	autoupdate.Init()

	err := settings.LoadSettings()
	if err != nil {
		slog.Error("failed to load settings", slog.Any("error", err))
		// Cannot use wails message dialogs here yet, because they expect a frontend to exist
		_ = dialog.Error("Failed to load settings: %s", err.Error())
		os.Exit(1)
	}

	if settings.Settings.CacheDir != "" {
		err = settings.ValidateCacheDir(settings.Settings.CacheDir)
		if err != nil {
			slog.Error("failed to set cache dir", slog.Any("error", err))
		} else {
			viper.Set("cache-dir", settings.Settings.CacheDir)
		}
	}

	b := bindings.MakeBindings()

	windowStartState := options.Normal
	if settings.Settings.Maximized {
		windowStartState = options.Maximised
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:            "SatisfactoryModManager",
		Frameless:        runtime.GOOS == "windows",
		Width:            settings.Settings.UnexpandedSize.Width,
		Height:           settings.Settings.UnexpandedSize.Height,
		MinWidth:         utils.UnexpandedMin.Width,
		MaxWidth:         utils.UnexpandedMax.Width,
		MinHeight:        utils.UnexpandedMin.Height,
		MaxHeight:        utils.UnexpandedMax.Height,
		WindowStartState: windowStartState,
		AssetServer:      &assetserver.Options{Assets: assets},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "SatisfactoryModManager_b04ab4c3-450f-48f4-ab14-af6d7adc5416",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				b.App.Show()
				backend.ProcessArguments(secondInstanceData.Args)
			},
		},
		OnStartup: func(ctx context.Context) {
			app.Context = ctx
			go websocket.ListenAndServeWebsocket()
			err := b.Startup(ctx)
			if err != nil {
				slog.Error("failed to create bindings", slog.Any("error", err))
				_ = dialog.Error("Failed to create bindings: %s", err.Error())
				os.Exit(1)
			}
		},
		OnDomReady: func(ctx context.Context) {
			backend.ProcessArguments(os.Args[1:])
			autoupdate.Updater.CheckInterval(5 * time.Minute)
		},
		OnShutdown: func(ctx context.Context) {
			b.Shutdown(ctx)
		},
		Bind: append(b.GetBindings(), autoupdate.Updater, settings.Settings),
		EnumBind: []interface{}{
			bindings.AllInstallTypes,
			bindings.AllBranches,
		},
		Logger: backend.WailsZeroLogLogger{},
	})

	if err != nil {
		slog.Error("failed to start application", slog.Any("error", err))
		_ = dialog.Error("Failed to start application: %s", err.Error())
	}

	err = autoupdate.Updater.OnExit()
	if err != nil {
		slog.Error("failed to apply update on exit", slog.Any("error", err))
		_ = dialog.Error("Failed to apply update on exit: %s", err.Error())
	}
}

func init() {
	// Pass build-time variables to viper
	if len(version) > 0 && version[0] == 'v' {
		version = version[1:]
	}
	viper.Set("version", version)
	viper.Set("commit", commit)
	viper.Set("date", date)
	viper.Set("update-mode", updateMode)

	// general config

	var baseLocalDir string

	switch runtime.GOOS {
	case "windows":
		baseLocalDir = os.Getenv("APPDATA")
	case "linux":
		baseLocalDir = path.Join(os.Getenv("HOME"), ".local", "share")
	default:
		slog.Error("unsupported OS", slog.String("os", runtime.GOOS))
		_ = dialog.Error("Unsupported OS: %s", runtime.GOOS)
		os.Exit(1)
	}

	viper.Set("base-local-dir", baseLocalDir)

	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	cacheDir := filepath.Clean(filepath.Join(baseCacheDir, "ficsit"))
	_ = utils.EnsureDirExists(cacheDir)
	viper.Set("cache-dir", cacheDir)
	viper.Set("default-cache-dir", cacheDir)

	smmCacheDir := filepath.Clean(filepath.Join(baseCacheDir, "SatisfactoryModManager"))
	_ = utils.EnsureDirExists(smmCacheDir)
	viper.Set("smm-cache-dir", smmCacheDir)

	localDir := filepath.Clean(filepath.Join(baseLocalDir, "ficsit"))
	_ = utils.EnsureDirExists(localDir)
	viper.Set("local-dir", localDir)

	smmLocalDir := filepath.Clean(filepath.Join(baseLocalDir, "SatisfactoryModManager"))
	_ = utils.EnsureDirExists(smmLocalDir)
	viper.Set("smm-local-dir", smmLocalDir)

	viper.Set("websocket-port", 33642)

	viper.Set("github-release-repo", "satisfactorymodding/SatisfactoryModManager")

	// ficsit-cli config
	viper.Set("profiles-file", "profiles.json")
	viper.Set("installations-file", "installations.json")
	viper.Set("api-base", "https://api.ficsit.app")
	viper.Set("graphql-api", "/v2/query")

	viper.Set("log", "info")
	viper.Set("log-file", filepath.Join(smmCacheDir, "logs", "SatisfactoryModManager.log"))

	viper.Set("concurrent-downloads", 5)

	level := slog.LevelInfo
	_ = level.UnmarshalText([]byte(viper.GetString("log")))

	handlers := make([]slog.Handler, 0)
	handlers = append(handlers, tint.NewHandler(os.Stdout, &tint.Options{
		Level:      level,
		AddSource:  true,
		TimeFormat: time.RFC3339,
	}))

	if viper.GetString("log-file") != "" {
		logFile := &lumberjack.Logger{
			Filename:   viper.GetString("log-file"),
			MaxSize:    10, // megabytes
			MaxBackups: 5,
			MaxAge:     30, // days
		}

		handlers = append(handlers, slog.NewJSONHandler(logFile, nil))
	}

	slog.SetDefault(slog.New(slogmulti.Fanout(handlers...)))
}
