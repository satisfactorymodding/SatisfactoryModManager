package heroic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/legendary"
)

func findInstallationsHeroic(snap bool, xdgConfigHomeEnv string, launcher string) ([]*common.Installation, []error) {
	legendaryDataPath, err := getHeroicLegendaryConfigPath(snap, xdgConfigHomeEnv)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get heroic legendary config paths: %w", err)}
	}

	return legendary.FindInstallationsIn(legendaryDataPath, launcher, common.MakeLauncherPlatform(common.NativePlatform(), nil))
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
