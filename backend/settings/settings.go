package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
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
	WindowPosition *utils.Position `json:"windowPosition"`
	Maximized      bool            `json:"maximized"`

	UnexpandedSize utils.Size `json:"unexpandedSize"`
	ExpandedSize   utils.Size `json:"expandedSize"`

	StartView View `json:"startView"`

	FavoriteMods []string        `json:"favoriteMods"`
	ModFilters   SavedModFilters `json:"modFilters"`

	QueueAutoStart      bool                `json:"queueAutoStart"`
	IgnoredUpdates      map[string][]string `json:"ignoredUpdates"`
	UpdateCheckMode     UpdateCheckMode     `json:"updateCheckMode"`
	ViewedAnnouncements []string            `json:"viewedAnnouncements"`

	Offline bool `json:"offline"`

	Konami       bool   `json:"konami"`
	LaunchButton string `json:"launchButton"`
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
			return errors.Wrap(err, "failed to stat settings file")
		}

		err = SaveSettings()
		if err != nil {
			return errors.Wrap(err, "failed to save default settings")
		}
	}

	settingsFile, err := os.ReadFile(filepath.Join(viper.GetString("smm-local-dir"), settingsFileName))
	if err != nil {
		return errors.Wrap(err, "failed to read settings")
	}

	if err := json.Unmarshal(settingsFile, &Settings); err != nil {
		// Settings file might be SMM2 settings, try to load those
		err = readSMM2Settings(settingsFile)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal settings")
		}
	}

	return nil
}

func SaveSettings() error {
	settingsFile, err := utils.JSONMarshal(Settings, 2)
	if err != nil {
		return errors.Wrap(err, "failed to marshal settings")
	}
	err = os.WriteFile(filepath.Join(viper.GetString("smm-local-dir"), settingsFileName), settingsFile, 0o755)
	if err != nil {
		return errors.Wrap(err, "failed to write settings")
	}

	return nil
}
