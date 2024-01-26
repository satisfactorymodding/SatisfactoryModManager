package main

import (
	"context"
	"embed"
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/lmittmann/tint"
	"github.com/pkg/errors"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/bindings"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/bindings/ficsitcli"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/projectfile"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/websocket"
)

//go:embed all:frontend/build
var assets embed.FS

//go:embed wails.json
var projectFile []byte

func loadProjectFile() error {
	err := json.Unmarshal(projectFile, &projectfile.ProjectFile)
	if err != nil {
		return errors.Wrap(err, "failed to load project file")
	}
	return nil
}

func main() {
	autoupdate.Init(autoupdate.Config{
		UpdateFoundCallback: func(latestVersion string, changelogs map[string]string) {
			bindings.BindingsInstance.Update.UpdateAvailable(latestVersion, changelogs)
		},
		DownloadProgressCallback: func(bytesDownloaded, totalBytes int64) {
			bindings.BindingsInstance.Update.UpdateDownloadProgress(bytesDownloaded, totalBytes)
		},
		UpdateReadyCallback: func() { bindings.BindingsInstance.Update.UpdateReady() },
	})

	err := settings.LoadSettings()
	if err != nil {
		slog.Error("failed to load settings", slog.Any("error", err))
		os.Exit(1)
	}

	if settings.Settings.CacheDir != "" {
		err = ficsitcli.ValidateCacheDir(settings.Settings.CacheDir)
		if err != nil {
			slog.Error("failed to set cache dir", slog.Any("error", err))
		} else {
			viper.Set("cache-dir", settings.Settings.CacheDir)
		}
	}

	b, err := bindings.MakeBindings()
	if err != nil {
		slog.Error("failed to create bindings", slog.Any("error", err))
		os.Exit(1)
	}

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
			go websocket.ListenAndServeWebsocket()
			b.Startup(ctx)
		},
		OnDomReady: func(ctx context.Context) {
			backend.ProcessArguments(os.Args[1:])
			autoupdate.CheckInterval(5 * time.Minute)
		},
		Bind: b.GetBindings(),
		EnumBind: []interface{}{
			bindings.AllInstallTypes,
			bindings.AllBranches,
		},
		Logger: backend.WailsZeroLogLogger{},
	})

	if err != nil {
		slog.Error("failed to start application", slog.Any("error", err))
	}

	err = autoupdate.OnExit(bindings.BindingsInstance.Update.Restart)
	if err != nil {
		slog.Error("failed to apply update on exit", slog.Any("error", err))
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
		slog.Error("unsupported OS", slog.String("os", runtime.GOOS))
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

	viper.Set("version", projectfile.ProjectFile.Info.ProductVersion)
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
