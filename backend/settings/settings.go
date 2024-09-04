package settings

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

type SavedModFilters struct {
	Order  string `json:"order"`
	Filter string `json:"filter"`
}

type View string

var (
	ViewCompact  View = "compact"
	ViewExpanded View = "expanded"
)

type UpdateCheckMode string

var (
	UpdateOnLaunch UpdateCheckMode = "launch"
	UpdateOnExit   UpdateCheckMode = "exit"
	UpdateAsk      UpdateCheckMode = "ask"
)

type settings struct {
	WindowPosition *utils.Position `json:"windowPosition,omitempty"`
	Maximized      bool            `json:"maximized,omitempty"`

	UnexpandedSize utils.Size `json:"unexpandedSize,omitempty"`
	ExpandedSize   utils.Size `json:"expandedSize,omitempty"`

	StartView View `json:"startView,omitempty"`

	FavoriteMods []string        `json:"favoriteMods,omitempty"`
	ModFilters   SavedModFilters `json:"modFilters,omitempty"`

	RemoteNames map[string]string `json:"remoteNames,omitempty"`

	QueueAutoStart      bool                `json:"queueAutoStart"`
	IgnoredUpdates      map[string][]string `json:"ignoredUpdates,omitempty"`
	UpdateCheckMode     UpdateCheckMode     `json:"updateCheckMode,omitempty"`
	ViewedAnnouncements []string            `json:"viewedAnnouncements,omitempty"`

	Offline bool `json:"offline,omitempty"`

	Language string `json:"language,omitempty"`

	Proxy string `json:"proxy,omitempty"`

	Konami       bool   `json:"konami,omitempty"`
	LaunchButton string `json:"launchButton,omitempty"`

	CacheDir string `json:"cacheDir,omitempty"`

	Debug bool `json:"debug,omitempty"`

	NewUserSetupComplete bool `json:"newUserSetupComplete,omitempty"`
}

var Settings = &settings{
	WindowPosition: nil,
	Maximized:      false,

	UnexpandedSize: utils.UnexpandedDefault,
	ExpandedSize:   utils.ExpandedDefault,

	StartView: ViewCompact,

	FavoriteMods: []string{},
	ModFilters: SavedModFilters{
		Order:  "last-updated",
		Filter: "compatible",
	},

	RemoteNames: map[string]string{},

	QueueAutoStart:      true,
	IgnoredUpdates:      map[string][]string{},
	UpdateCheckMode:     UpdateOnLaunch,
	ViewedAnnouncements: []string{},

	Offline: false,

	Konami:       false,
	LaunchButton: "normal",

	Debug: false,

	NewUserSetupComplete: false,
}

func (s *settings) GetNewUserSetupComplete() bool {
	return s.Debug
}

func (s *settings) SetNewUserSetupComplete(value bool) {
	slog.Info("changing NewUserSetupComplete state", slog.Bool("value", value))
	s.NewUserSetupComplete = value
	_ = SaveSettings()
}

func (s *settings) FavoriteMod(modReference string) (bool, error) {
	idx := -1
	for i, mod := range s.FavoriteMods {
		if mod == modReference {
			idx = i
			break
		}
	}
	if idx != -1 {
		return false, nil
	}
	s.FavoriteMods = append(s.FavoriteMods, modReference)
	err := SaveSettings()
	if err != nil {
		return false, err
	}
	s.emitFavoriteMods()
	return true, nil
}

func (s *settings) UnFavoriteMod(modReference string) bool {
	idx := -1
	for i, mod := range s.FavoriteMods {
		if mod == modReference {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}
	s.FavoriteMods = append(s.FavoriteMods[:idx], s.FavoriteMods[idx+1:]...)
	_ = SaveSettings()
	s.emitFavoriteMods()
	return true
}

func (s *settings) GetFavoriteMods() []string {
	return s.FavoriteMods
}

func (s *settings) GetModFiltersOrder() string {
	return s.ModFilters.Order
}

func (s *settings) GetModFiltersFilter() string {
	return s.ModFilters.Filter
}

func (s *settings) SetModFiltersOrder(order string) {
	s.ModFilters.Order = order
	_ = SaveSettings()
}

func (s *settings) SetModFiltersFilter(filter string) {
	s.ModFilters.Filter = filter
	_ = SaveSettings()
}

func (s *settings) emitFavoriteMods() {
	wailsRuntime.EventsEmit(common.AppContext, "favoriteMods", s.FavoriteMods)
}

func (s *settings) GetStartView() View {
	return s.StartView
}

func (s *settings) SetStartView(view View) {
	s.StartView = view
	_ = SaveSettings()
}

func (s *settings) GetKonami() bool {
	return s.Konami
}

func (s *settings) SetKonami(value bool) {
	s.Konami = value
	_ = SaveSettings()
}

func (s *settings) GetLaunchButton() string {
	return s.LaunchButton
}

func (s *settings) SetLaunchButton(value string) {
	s.LaunchButton = value
	_ = SaveSettings()
}

func (s *settings) GetQueueAutoStart() bool {
	return s.QueueAutoStart
}

func (s *settings) SetQueueAutoStart(value bool) {
	s.QueueAutoStart = value
	_ = SaveSettings()
}

func (s *settings) GetIgnoredUpdates() map[string][]string {
	return s.IgnoredUpdates
}

func (s *settings) SetUpdateIgnore(modReference string, version string) {
	s.IgnoredUpdates[modReference] = append(s.IgnoredUpdates[modReference], version)
	_ = SaveSettings()
	wailsRuntime.EventsEmit(common.AppContext, "ignoredUpdates", s.IgnoredUpdates)
}

func (s *settings) SetUpdateUnignore(modReference string, version string) {
	versions := s.IgnoredUpdates[modReference]
	idx := -1
	for i, v := range versions {
		if v == version {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}
	s.IgnoredUpdates[modReference] = append(versions[:idx], versions[idx+1:]...)
	_ = SaveSettings()
	wailsRuntime.EventsEmit(common.AppContext, "ignoredUpdates", s.IgnoredUpdates)
}

func (s *settings) GetUpdateCheckMode() UpdateCheckMode {
	return s.UpdateCheckMode
}

func (s *settings) SetUpdateCheckMode(value UpdateCheckMode) {
	s.UpdateCheckMode = value
	_ = SaveSettings()
}

func (s *settings) GetViewedAnnouncements() []string {
	return s.ViewedAnnouncements
}

func (s *settings) SetAnnouncementViewed(announcement string) {
	found := false
	for _, viewed := range s.ViewedAnnouncements {
		if viewed == announcement {
			found = true
			break
		}
	}
	if found {
		return
	}
	s.ViewedAnnouncements = append(s.ViewedAnnouncements, announcement)
	_ = SaveSettings()
	wailsRuntime.EventsEmit(common.AppContext, "viewedAnnouncements", s.ViewedAnnouncements)
}

func (s *settings) GetLanguage() string {
	return s.Language
}

func (s *settings) SetLanguage(value string) {
	s.Language = value
	_ = SaveSettings()
}

func (s *settings) GetDebug() bool {
	return s.Debug
}

func (s *settings) SetDebug(value bool) {
	slog.Info("changing debug mode state", slog.Bool("value", value))
	s.Debug = value
	_ = SaveSettings()
}

func (s *settings) GetProxy() string {
	return s.Proxy
}

func (s *settings) SetProxy(value string) {
	s.Proxy = value
	_ = SaveSettings()
}

func (s *settings) SetCacheDir(dir string) error {
	realDir := dir
	if dir == "" {
		realDir = viper.GetString("default-cache-dir")
	}
	err := moveCacheDir(realDir)
	if err != nil {
		slog.Error("failed to set cache dir", slog.Any("error", err))
		return err
	}
	s.CacheDir = dir
	_ = SaveSettings()
	wailsRuntime.EventsEmit(common.AppContext, "cacheDir", s.GetCacheDir())
	return nil
}

func (s *settings) GetCacheDir() string {
	return viper.GetString("cache-dir")
}

func ValidateCacheDir(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat %s: %w", dir, err)
		}
	} else {
		if !stat.IsDir() {
			return fmt.Errorf("%s is not a directory", dir)
		}
	}
	return nil
}

func moveCacheDir(newDir string) error {
	if newDir == viper.GetString("cache-dir") {
		return nil
	}

	err := ValidateCacheDir(newDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(newDir, 0o755)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("failed to create %s: %w", newDir, err)
		}
	}

	items, err := os.ReadDir(newDir)
	if err != nil {
		return fmt.Errorf("failed to check if directory %s is empty: %w", newDir, err)
	}
	if len(items) > 0 {
		return fmt.Errorf("directory %s is not empty", newDir)
	}

	oldCacheDir := viper.GetString("cache-dir")
	// Move contents of oldCacheDir to dir
	if oldCacheDir != "" && oldCacheDir != newDir {
		err := moveCacheData(oldCacheDir, newDir)
		if err != nil {
			return err
		}
	}

	viper.Set("cache-dir", newDir)
	return nil
}

func moveCacheData(oldCacheDir, newDir string) error {
	oldStat, err := os.Stat(oldCacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Nothing to move
			return nil
		}
		return fmt.Errorf("failed to stat %s: %w", oldCacheDir, err)
	}
	if !oldStat.IsDir() {
		return fmt.Errorf("%s is not a directory", oldCacheDir)
	}

	// Perform the move atomically
	copySuccess, err := utils.MoveRecursive(oldCacheDir, newDir)
	if err != nil {
		if !copySuccess {
			return err
		}
		slog.Error("failed to move cache dir", slog.Any("error", err))
	}

	return nil
}

var settingsFileName = "settings.json"

func LoadSettings() error {
	settingsFilePath := filepath.Join(viper.GetString("smm-local-dir"), settingsFileName)

	_, err := os.Stat(settingsFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat settings file: %w", err)
		}

		err = SaveSettings()
		if err != nil {
			return fmt.Errorf("failed to save default settings: %w", err)
		}
	}

	settingsFile, err := os.ReadFile(filepath.Join(viper.GetString("smm-local-dir"), settingsFileName))
	if err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	if err := json.Unmarshal(settingsFile, Settings); err != nil {
		// Settings file might be SMM2 settings, try to load those
		err = readSMM2Settings(settingsFile)
		if err != nil {
			return fmt.Errorf("failed to unmarshal settings: %w", err)
		}
	}

	return nil
}

func SaveSettings() error {
	settingsFile, err := utils.JSONMarshal(Settings, 2)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}
	err = os.WriteFile(filepath.Join(viper.GetString("smm-local-dir"), settingsFileName), settingsFile, 0o755)
	if err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}
