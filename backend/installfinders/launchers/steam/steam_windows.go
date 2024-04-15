package steam

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows/registry"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Steam", func() ([]*common.Installation, []error) {
		steamPath, err := getSteamPath()
		if err != nil {
			return nil, []error{err}
		}

		return FindInstallationsSteam(
			steamPath,
			"Steam",
			func(steamApp string) []string {
				return []string{
					"cmd",
					"/C",
					"start",
					"",
					steamApp,
				}
			},
			nil,
		)
	})
}

func getSteamPath() (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("failed to open Steam registry key: %w", err)
	}
	defer key.Close()

	steamExePath, _, err := key.GetStringValue("SteamExe")
	if err != nil {
		steamExePath = `C:\Program Files (x86)\Steam\steam.exe`
	}

	return filepath.Dir(steamExePath), nil
}
