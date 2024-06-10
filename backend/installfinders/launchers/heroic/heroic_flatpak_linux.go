package heroic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Heroic-flatpak", func() ([]*common.Installation, []error) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
		}
		flatpakXdgConfigHome := filepath.Join(homeDir, ".var", "app", "com.heroicgameslauncher.hgl", "config")

		return findInstallationsHeroic(false, flatpakXdgConfigHome, "Heroic")
	})
}
