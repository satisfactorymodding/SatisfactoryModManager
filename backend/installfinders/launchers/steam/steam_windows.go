package steam

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows/registry"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func FindInstallations() ([]*common.Installation, []error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to open Steam registry key: %w", err)}
	}
	defer key.Close()

	steamExePath, _, err := key.GetStringValue("SteamExe")
	if err != nil {
		steamExePath = `C:\Program Files (x86)\Steam\steam.exe`
	}

	steamPath := filepath.Dir(steamExePath)
	return findInstallationsSteam(
		steamPath,
		"Steam",
		[]string{
			"cmd",
			"/C",
			"start",
			"",
		},
	)
}
