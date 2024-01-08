package settings

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

type smm2Settings struct {
	WindowLocation *struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"windowLocation"`
	NormalSize *struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"normalSize"`
	ExpandedSize *struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"expandedSize"`
	FavoriteMods *[]string `json:"favoriteMods"`
	Filters      *struct {
		ModFilters string `json:"modFilters"`
		SortBy     string `json:"sortBy"`
	} `json:"filters"`
	IgnoredUpdates *[]struct {
		Item    string `json:"item"`
		Version string `json:"version"`
	} `json:"ignoredUpdates"`
	Maximized              *bool              `json:"maximized"`
	SelectedInstall        *string            `json:"selectedInstall"`
	DebugMode              *bool              `json:"debugMode"`
	SelectedProfile        *map[string]string `json:"selectedProfile"`
	UpdateCheckMode        *UpdateCheckMode   `json:"updateCheckMode"`
	ModsEnabled            *map[string]bool   `json:"modsEnabled"`
	Konami                 *bool              `json:"konami"`
	LaunchButton           *bool              `json:"launchButton"`
	ExpandModInfoOnStart   *bool              `json:"expandModInfoOnStart"`
	LaunchCat              *bool              `json:"launchCat"`
	ViewedAnnouncements    *[]string          `json:"viewedAnnouncements"`
	DisableDownloadTimeout *bool              `json:"disableDownloadTimeout"`
}

func readSMM2Settings(data []byte) error {
	s := smm2Settings{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal SMM2 settings")
	}

	if s.WindowLocation != nil {
		Settings.WindowPosition = &utils.Position{
			X: s.WindowLocation.X,
			Y: s.WindowLocation.Y,
		}
	}

	if s.NormalSize != nil {
		Settings.UnexpandedSize.Width = s.NormalSize.Width
		Settings.UnexpandedSize.Height = s.NormalSize.Height
	}

	if s.ExpandedSize != nil {
		Settings.ExpandedSize.Width = s.ExpandedSize.Width
		Settings.ExpandedSize.Height = s.ExpandedSize.Height
	}

	if s.FavoriteMods != nil {
		Settings.FavoriteMods = *s.FavoriteMods
	}

	if s.Filters != nil {
		Settings.ModFilters.Filter = s.Filters.ModFilters
		Settings.ModFilters.Order = s.Filters.SortBy
	}

	if s.IgnoredUpdates != nil {
		for _, ignoredUpdate := range *s.IgnoredUpdates {
			Settings.IgnoredUpdates[ignoredUpdate.Item] = append(Settings.IgnoredUpdates[ignoredUpdate.Item], ignoredUpdate.Version)
		}
	}

	if s.Maximized != nil {
		Settings.Maximized = *s.Maximized
	}

	// Ignore selected install, profile, and mods enabled
	// They are stored in ficsit-cli, but that gets initialized later
	// They are not critical anyway

	// Ignore DebugMode, it's not used anymore

	if s.UpdateCheckMode != nil {
		Settings.UpdateCheckMode = *s.UpdateCheckMode
	}

	if s.Konami != nil {
		Settings.Konami = *s.Konami
	}

	if s.LaunchButton != nil {
		if *s.LaunchButton {
			Settings.LaunchButton = "button"
		}
	}

	if s.LaunchCat != nil {
		if *s.LaunchCat {
			Settings.LaunchButton = "cat"
		}
	}

	if s.ExpandModInfoOnStart != nil {
		if *s.ExpandModInfoOnStart {
			Settings.StartView = ViewExpanded
		}
	}

	if s.ViewedAnnouncements != nil {
		Settings.ViewedAnnouncements = *s.ViewedAnnouncements
	}

	// Ignore DisableDownloadTimeout, it's not used anymore

	return nil
}
