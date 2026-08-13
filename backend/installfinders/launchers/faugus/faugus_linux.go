package faugus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/epic"
)

type FaugusGame struct {
	Gameid          string `json:"gameid"`
	Title           string `json:"title"`
	Path            string `json:"path"`
	Prefix          string `json:"prefix"`
	LaunchArguments string `json:"launch_arguments"`
	GameArguments   string `json:"game_arguments"`
}

var faugusGamesRelativePath = filepath.Join("faugus-launcher", "games.json")

func init() {
	launchers.Add("Faugus", func() ([]*common.Installation, []error) {
		dataDir, ok := os.LookupEnv("XDG_DATA_HOME")
		if !ok {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
			}
			dataDir = filepath.Join(homeDir, ".local", "share")
		}

		return findInstallationsFaugus(dataDir, "Faugus", []string{"faugus-launcher"})
	})
	launchers.Add("Faugus-flatpak", func() ([]*common.Installation, []error) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
		}
		flatpakXdgDataHome := filepath.Join(homeDir, ".var", "app", "io.github.Faugus.faugus-launcher", "data")

		return findInstallationsFaugus(flatpakXdgDataHome, "Faugus", []string{"flatpak", "run", "io.github.Faugus.faugus-launcher"})
	})
}

func findInstallationsFaugus(xdgDataHome string, launcher string, launchCommand []string) ([]*common.Installation, []error) {
	faugusGamesJSONPath := filepath.Join(xdgDataHome, faugusGamesRelativePath)

	faugusGamesJSON, err := os.ReadFile(faugusGamesJSONPath)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to read faugus games.json: %w", err)}
	}

	var faugusGames []FaugusGame
	err = json.Unmarshal(faugusGamesJSON, &faugusGames)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to parse faugus games.json: %w", err)}
	}

	var epicGamesGame *FaugusGame
	for _, game := range faugusGames {
		if game.Gameid == "epic-games" {
			epicGamesGame = &game
			break
		}
	}

	if epicGamesGame == nil {
		return nil, nil
	}

	return epic.FindInstallationsWine(epicGamesGame.Prefix, launcher, append(launchCommand, "--game", epicGamesGame.Gameid))
}
