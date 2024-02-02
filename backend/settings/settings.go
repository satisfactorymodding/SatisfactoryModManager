package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

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

	QueueAutoStart      bool                `json:"queueAutoStart,omitempty"`
	IgnoredUpdates      map[string][]string `json:"ignoredUpdates,omitempty"`
	UpdateCheckMode     UpdateCheckMode     `json:"updateCheckMode,omitempty"`
	ViewedAnnouncements []string            `json:"viewedAnnouncements,omitempty"`

	Offline bool `json:"offline,omitempty"`

	Konami       bool   `json:"konami,omitempty"`
	LaunchButton string `json:"launchButton,omitempty"`

	CacheDir string `json:"cacheDir,omitempty"`
}

var Settings = settings{
	WindowPosition: nil,
	Maximized:      false,

	UnexpandedSize: utils.UnexpandedDefault,
	ExpandedSize:   utils.ExpandedDefault,

	StartView: ViewCompact,

	FavoriteMods: []string{},
	ModFilters: SavedModFilters{
		Order:  "Last updated",
		Filter: "Compatible",
	},

	QueueAutoStart:      true,
	IgnoredUpdates:      map[string][]string{},
	UpdateCheckMode:     UpdateOnLaunch,
	ViewedAnnouncements: []string{},

	Offline: false,

	Konami:       false,
	LaunchButton: "normal",
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

	if err := json.Unmarshal(settingsFile, &Settings); err != nil {
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
