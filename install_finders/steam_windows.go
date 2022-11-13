package install_finders

import (
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

func FindInstallationsWindowsSteam() ([]*Installation, []error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return nil, []error{errors.Wrap(err, "Failed to open Steam registry key")}
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
