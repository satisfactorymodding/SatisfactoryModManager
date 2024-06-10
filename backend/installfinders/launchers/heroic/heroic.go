package heroic

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/legendary"
)

func findInstallationsHeroic(snap bool, xdgConfigHomeEnv string, launcher string) ([]*common.Installation, []error) {
	legendaryDataPath, err := getHeroicLegendaryConfigPath(snap, xdgConfigHomeEnv)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get heroic legendary config paths: %w", err)}
	}

	knownPrefixes, err := getHeroicKnownWinePrefixes(xdgConfigHomeEnv)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get heroic known wine prefixes: %w", err)}
	}

	return legendary.FindInstallationsIn(legendaryDataPath, launcher, knownPrefixes, common.MakeLauncherPlatform(common.NativePlatform(), nil))
}

func getHeroicKnownWinePrefixes(xdgConfigHomeEnv string) (map[string]string, error) {
	configPath := xdgConfigHomeEnv
	if configPath == "" {
		var err error
		configPath, err = os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user config dir: %w", err)
		}
	}

	heroicGamesConfigPath := filepath.Join(configPath, "heroic", "GamesConfig")

	knownPrefixes := make(map[string]string)

	items, err := os.ReadDir(heroicGamesConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return knownPrefixes, nil
		}
		return nil, fmt.Errorf("failed to read heroic games config path: %w", err)
	}

	for _, item := range items {
		if item.IsDir() || !strings.HasSuffix(item.Name(), ".json") {
			continue
		}

		gameID := strings.TrimSuffix(item.Name(), ".json")

		bytes, err := os.ReadFile(filepath.Join(heroicGamesConfigPath, item.Name()))
		if err != nil {
			slog.Warn("failed to read heroic game config", slog.String("path", item.Name()), slog.Any("error", err))
			continue
		}

		var game map[string]interface{}
		if err = json.Unmarshal(bytes, &game); err != nil {
			slog.Warn("failed to parse heroic game config", slog.String("path", item.Name()), slog.Any("error", err))
			continue
		}

		if gameEntry := game[gameID]; gameEntry != nil {
			gameData := gameEntry.(map[string]interface{})

			prefix, ok := gameData["winePrefix"].(string)
			if !ok {
				continue
			}

			knownPrefixes[gameID] = prefix
		}
	}

	return knownPrefixes, nil
}

func getHeroicLegendaryConfigPath(snap bool, xdgConfigHomeEnv string) (string, error) {
	// Allow passing xdgConfigHomeEnv for flatpak and snap support

	// Should be kept in sync with
	// https://github.com/Heroic-Games-Launcher/HeroicGamesLauncher/blob/main/src/backend/constants.ts#L56

	if snap {
		if xdgConfigHomeEnv == "" {
			return "", fmt.Errorf("creating path for heroic snap but XDG_CONFIG_HOME not set")
		}
		return filepath.Join(xdgConfigHomeEnv, "legendary"), nil
	}

	configPath := xdgConfigHomeEnv
	if configPath == "" {
		var err error
		configPath, err = os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config dir: %w", err)
		}
	}

	return filepath.Join(configPath, "heroic", "legendaryConfig", "legendary"), nil
}
