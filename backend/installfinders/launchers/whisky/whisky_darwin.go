package whisky

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"howett.net/plist"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/steam"
)

var (
	whiskyConfigRelativePath         = filepath.Join("Library", "Preferences", "com.isaacmarovitz.Whisky.plist")
	whiskyDefaultBottlesRelativePath = filepath.Join("Library", "Containers", "com.isaacmarovitz.Whisky", "Bottles")
	whiskyBottleVMRelativePath       = "BottleVM.plist"
	whiskySteamPath                  = filepath.Join("c:", "Program Files (x86)", "Steam") // Will get run through processPath, so it will be added to the dosdevices path
)

func init() {
	launchers.Add("Whisky", whisky)
}

func whisky() ([]*common.Installation, []error) {
	bottlesPath, err := getWhiskyBottlesPath()
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get Whisky bottles path: %w", err)}
	}

	if _, err := os.Stat(bottlesPath); os.IsNotExist(err) {
		return nil, []error{fmt.Errorf("whisky not installed")}
	}

	bottles, err := os.ReadDir(bottlesPath)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to list Whisky bottles: %w", err)}
	}

	bottlesToCheck := make([]string, 0, len(bottles))

	for _, bottle := range bottles {
		if !bottle.IsDir() {
			continue
		}
		bottleRoot := filepath.Join(bottlesPath, bottle.Name())
		bottlesToCheck = append(bottlesToCheck, bottleRoot)
	}

	bottleVMBottles, err := getBottlesFromBottlesVM(bottlesPath)
	if err != nil {
		slog.Error("failed to get list of additional whisky bottles", slog.Any("error", err))
	}
	bottlesToCheck = append(bottlesToCheck, bottleVMBottles...)

	installations := make([]*common.Installation, 0)
	errors := make([]error, 0)
	for _, bottleRoot := range bottlesToCheck {
		processPath := common.WinePathProcessor(bottleRoot)

		if _, err := os.Stat(processPath(whiskySteamPath)); os.IsNotExist(err) {
			slog.Debug("Skipping bottle without Steam", slog.String("bottle", bottleRoot))
			continue
		}
		bottleInstalls, bottleErrs := steam.FindInstallationsSteam(
			whiskySteamPath,
			"Whisky",
			func(steamApp string) []string {
				return nil
			},
			processPath,
		)
		installations = append(installations, bottleInstalls...)
		if bottleErrs != nil {
			errors = append(errors, bottleErrs...)
		}
	}

	return installations, errors
}

func getWhiskyBottlesPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}

	defaultBottlesPath := filepath.Join(homeDir, whiskyDefaultBottlesRelativePath)

	var bottlesPath string

	configPath := filepath.Join(homeDir, whiskyConfigRelativePath)
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("whisky config file missing")
		} else {
			slog.Error("failed to read Whisky config file", slog.Any("error", err))
		}
	} else {
		var config whiskyPlist
		_, err := plist.Unmarshal(configBytes, &config)
		if err != nil {
			slog.Error("failed to parse Whisky config file", slog.Any("error", err))
		} else {
			bottlesPath = config.DefaultBottleLocation
		}
	}

	if bottlesPath == "" {
		bottlesPath = defaultBottlesPath
	}

	return bottlesPath, nil
}

func getBottlesFromBottlesVM(bottlesPath string) ([]string, error) {
	bottleVMPath := filepath.Join(bottlesPath, whiskyBottleVMRelativePath)
	bottleVMBytes, err := os.ReadFile(bottleVMPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read BottleVM.plist: %w", err)
	}

	var bottleVM bottleVMPlist
	_, err = plist.Unmarshal(bottleVMBytes, &bottleVM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse BottleVM.plist: %w", err)
	}

	bottles := make([]string, 0, len(bottleVM.Paths))
	for _, path := range bottleVM.Paths {
		parsed, err := url.Parse(path.Relative)
		if err != nil {
			return nil, fmt.Errorf("failed to parse BottleVM path: %w", err)
		}
		// Even through the name is "relative", it's actually an absolute path stored as a file:// URL
		bottles = append(bottles, parsed.Path)
	}

	return bottles, nil
}
