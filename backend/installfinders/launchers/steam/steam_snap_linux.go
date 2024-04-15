package steam

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Steam-snap", func() ([]*common.Installation, []error) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
		}
		steamPath := filepath.Join(homeDir, "snap", "steam", "common", ".local", "share", "Steam")
		if _, err := os.Stat(steamPath); os.IsNotExist(err) {
			return nil, []error{fmt.Errorf("steam-snap not installed")}
		}
		return FindInstallationsSteam(
			steamPath,
			"Steam",
			func(steamApp string) []string {
				return []string{
					"snap",
					"run",
					"steam",
					steamApp,
				}
			},
			nil,
		)
	})
}
