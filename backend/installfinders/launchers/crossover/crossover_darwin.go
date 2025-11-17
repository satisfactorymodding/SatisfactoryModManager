package crossover

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"howett.net/plist"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/epic"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/steam"
)

var (
	crossoverConfigRelativePath         = filepath.Join("Library", "Preferences", "com.codeweavers.CrossOver.plist")
	crossoverDefaultBottlesRelativePath = filepath.Join("Library", "Application Support", "Crossover", "Bottles")
)

func init() {
	launchers.Add("CrossOver", crossover)
}

func crossover() ([]*common.Installation, []error) {
	bottlesPath, err := getCrossoverBottlesPath()
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get CrossOver bottles path: %w", err)}
	}

	if _, err := os.Stat(bottlesPath); os.IsNotExist(err) {
		return nil, []error{fmt.Errorf("crossover not installed")}
	}

	bottles, err := os.ReadDir(bottlesPath)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to list Crossover bottles: %w", err)}
	}

	installations := make([]*common.Installation, 0)
	errors := make([]error, 0)
	for _, bottle := range bottles {
		if !bottle.IsDir() {
			continue
		}
		bottleRoot := filepath.Join(bottlesPath, bottle.Name())
		search := func(finder func(winePrefix string, launcher string, launchPath []string) ([]*common.Installation, []error)) {
			bottleInstalls, bottleErrs := finder(bottleRoot, "CrossOver", nil)
			installations = append(installations, bottleInstalls...)
			if bottleErrs != nil {
				errors = append(errors, bottleErrs...)
			}
		}
		search(steam.FindInstallationsWine)
		search(epic.FindInstallationsWine)
	}

	return installations, errors
}

func getCrossoverBottlesPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}

	defaultBottlesPath := filepath.Join(homeDir, crossoverDefaultBottlesRelativePath)

	var bottlesPath string

	configPath := filepath.Join(homeDir, crossoverConfigRelativePath)
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("CrossOver config file missing")
		} else {
			slog.Error("failed to read CrossOver config file", slog.Any("error", err))
		}
	} else {
		var config crossoverPlist
		_, err := plist.Unmarshal(configBytes, &config)
		if err != nil {
			slog.Error("failed to parse CrossOver config file", slog.Any("error", err))
		} else {
			bottlesPath = config.BottleDir
		}
	}

	if bottlesPath == "" {
		bottlesPath = defaultBottlesPath
	}

	return bottlesPath, nil
}
