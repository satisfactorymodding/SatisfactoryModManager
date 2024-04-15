package steam

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Steam-flatpak", func() ([]*common.Installation, []error) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
		}

		steamPath := filepath.Join(homeDir, ".var", "app", "com.valvesoftware.Steam", ".steam", "steam")
		if _, err := os.Stat(steamPath); os.IsNotExist(err) {
			return nil, []error{fmt.Errorf("steam-flatpak not installed")}
		}
		return FindInstallationsSteam(
			steamPath,
			"Steam",
			func(steamApp string) []string {
				return []string{
					"flatpak",
					"run",
					"com.valvesoftware.Steam",
					steamApp,
				}
			},
			nil,
		)
	})
}
