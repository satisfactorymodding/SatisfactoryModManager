//go:build unix

package steam

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Steam", func() ([]*common.Installation, []error) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
		}

		steamPath := filepath.Join(homeDir, ".steam", "steam")
		if _, err := os.Stat(steamPath); os.IsNotExist(err) {
			return nil, []error{fmt.Errorf("steam not installed")}
		}

		return FindInstallationsSteam(
			steamPath,
			"Steam",
			common.MakeLauncherPlatform(
				common.NativePlatform(),
				func(steamApp string) []string {
					return []string{
						"steam",
						steamApp,
					}
				}),
		)
	})
}
