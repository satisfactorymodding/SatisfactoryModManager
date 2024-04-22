package main

import (
	"context"
	"embed"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
	"github.com/tawesoft/golib/v2/dialog"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/app"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate"
	appCommon "github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/ficsitcli"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/logging"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/wailsextras"
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
	logging.Init()

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

	if settings.Settings.Proxy != "" {
		// webkit honors these env vars, even if they are an empty string,
		// so we must ensure they are valid
		_, err := url.Parse(settings.Settings.Proxy)
		if err != nil {
			slog.Error("skipping setting proxy, invalid URL", slog.Any("error", err))
		} else {
			err = os.Setenv("HTTP_PROXY", settings.Settings.Proxy)
			if err != nil {
				slog.Error("failed to set HTTP_PROXY", slog.Any("error", err))
			}
			err = os.Setenv("HTTPS_PROXY", settings.Settings.Proxy)
			if err != nil {
				slog.Error("failed to set HTTPS_PROXY", slog.Any("error", err))
			}
		}
	}

	err = ficsitcli.Init()
	if err != nil {
		slog.Error("failed to initialize ficsit-cli", slog.Any("error", err))
		_ = dialog.Error("Failed to initialize ficsit-cli: %s", err.Error())
		os.Exit(1)
	}

	windowStartState := options.Normal
	if settings.Settings.Maximized {
		windowStartState = options.Maximised
	}

	if len(os.Args) > 1 && os.Args[1] == "wipe-mods" {
		includeRemote := len(os.Args) > 2 && os.Args[2] == "remote"
		err := ficsitcli.FicsitCLI.WipeMods(includeRemote)
		if err != nil {
			slog.Error("failed to wipe mods", slog.Any("error", err))
			_ = dialog.Error("Failed to wipe mods: %s", err.Error())
			os.Exit(1)
		}
		return
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
				app.App.Show()
				backend.ProcessArguments(secondInstanceData.Args)
			},
		},
		OnStartup: func(ctx context.Context) {
			appCommon.AppContext = ctx

			// Wails doesn't support setting the window position on init, so we do it here
			if settings.Settings.WindowPosition != nil {
				wailsextras.WindowSetPosition(ctx, settings.Settings.WindowPosition.X, settings.Settings.WindowPosition.Y)
			}

			app.App.WatchWindow() //nolint:contextcheck
			go websocket.ListenAndServeWebsocket()

			ficsitcli.FicsitCLI.StartGameRunningWatcher() //nolint:contextcheck
		},
		OnDomReady: func(ctx context.Context) {
			backend.ProcessArguments(os.Args[1:]) //nolint:contextcheck
			autoupdate.Updater.CheckInterval(5 * time.Minute)
		},
		OnShutdown: func(ctx context.Context) {
			app.App.StopWindowWatcher()
		},
		Bind: []interface{}{
			app.App,
			ficsitcli.FicsitCLI,
			autoupdate.Updater,
			settings.Settings,
			ficsitcli.ServerPicker,
		},
		EnumBind: []interface{}{
			common.AllInstallTypes,
			common.AllBranches,
			common.AllLocationTypes,
			ficsitcli.AllInstallationStates,
			ficsitcli.AllActionTypes,
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

	if app.App.Restart && !autoupdate.Updater.HasRestarted() {
		err := utils.Restart()
		if err != nil {
			slog.Error("failed to restart", slog.Any("error", err))
			_ = dialog.Error("Failed to restart: %s", err.Error())
		}
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

	var baseLocalDir string

	switch runtime.GOOS {
	case "linux":
		baseLocalDir = filepath.Join(os.Getenv("HOME"), ".local", "share")
	default:
		var err error
		baseLocalDir, err = os.UserConfigDir()
		if err != nil {
			slog.Error("failed to get config dir", slog.Any("error", err))
			_ = dialog.Error("Failed to get config dir", err.Error())
			os.Exit(1)
		}
	}

	viper.Set("base-local-dir", baseLocalDir)

	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	// ficsit-cli config

	viper.Set("profiles-file", "profiles.json")
	viper.Set("installations-file", "installations.json")
	viper.Set("api-base", "https://api.ficsit.app")
	viper.Set("graphql-api", "/v2/query")
	viper.Set("concurrent-downloads", 5)

	cacheDir := filepath.Clean(filepath.Join(baseCacheDir, "ficsit"))
	_ = utils.EnsureDirExists(cacheDir)
	viper.Set("cache-dir", cacheDir)

	localDir := filepath.Clean(filepath.Join(baseLocalDir, "ficsit"))
	_ = utils.EnsureDirExists(localDir)
	viper.Set("local-dir", localDir)

	// SMM config

	smmCacheDir := filepath.Clean(filepath.Join(baseCacheDir, "SatisfactoryModManager"))
	_ = utils.EnsureDirExists(smmCacheDir)
	viper.Set("smm-cache-dir", smmCacheDir)

	smmLocalDir := filepath.Clean(filepath.Join(baseLocalDir, "SatisfactoryModManager"))
	_ = utils.EnsureDirExists(smmLocalDir)
	viper.Set("smm-local-dir", smmLocalDir)

	viper.Set("default-cache-dir", cacheDir)

	viper.Set("websocket-port", 33642)

	viper.Set("github-release-repo", "satisfactorymodding/SatisfactoryModManager")

	// logging

	viper.Set("log-file", filepath.Join(smmCacheDir, "logs", "SatisfactoryModManager.log"))
}
