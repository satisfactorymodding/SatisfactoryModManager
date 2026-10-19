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

		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			dataHome = filepath.Join(homeDir, ".local", "share")
		}

		steamPaths := []string{
			filepath.Join(homeDir, ".steam", "steam"),
			filepath.Join(dataHome, "Steam"),
			filepath.Join(dataHome, "SteamBeta"),
		}

		steamHome := os.Getenv("SMM_STEAM_HOME")
		if steamHome != "" {
			steamPaths = append(steamPaths, filepath.Join(steamHome, ".steam", "steam"))
		}

		realSteamPaths := make([]string, 0, len(steamPaths))
		for _, path := range steamPaths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}
			for _, existingPath := range realSteamPaths {
				if common.OsPathEqual(existingPath, path) {
					continue
				}
			}
			realSteamPaths = append(realSteamPaths, path)
		}

		if len(realSteamPaths) == 0 {
			return nil, []error{fmt.Errorf("steam not installed")}
		}

		var installations []*common.Installation
		var errors []error

		for _, steamPath := range steamPaths {
			i, e := FindInstallationsSteam(
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
			installations = append(installations, i...)
			errors = append(errors, e...)
		}

		return installations, errors
	})
}
