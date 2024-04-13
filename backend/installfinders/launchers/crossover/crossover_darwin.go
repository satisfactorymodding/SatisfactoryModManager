package crossover

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/steam"
)

var (
	crossoverBottlesRelativePath = filepath.Join("Library", "Application Support", "Crossover", "Bottles")
	crossoverSteamPath           = filepath.Join("c:", "Program Files (x86)", "Steam") // Will get run through processPath, so it will be added to the dosdevices path
)

func FindInstallations() ([]*common.Installation, []error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
	}
	bottlesPath := filepath.Join(homeDir, crossoverBottlesRelativePath)
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

		processPath := func(path string) string {
			return filepath.Join(bottleRoot, "dosdevices", strings.ToLower(path[0:1])+strings.ReplaceAll(path[1:], "\\", "/"))
		}

		if _, err := os.Stat(processPath(crossoverSteamPath)); os.IsNotExist(err) {
			slog.Debug("Skipping bottle without Steam", slog.String("bottle", bottle.Name()))
			continue
		}
		bottleInstalls, bottleErrs := steam.FindInstallationsSteam(
			crossoverSteamPath,
			"CrossOver",
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
