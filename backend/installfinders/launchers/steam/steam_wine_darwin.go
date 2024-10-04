package steam

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

// Will get run through processPath, so it will be added to the dosdevices path
// Theoretically this could be configured on custom wine prefixes (so would require parsing the wine registry),
// but the supported launchers don't, and have no reason to
var steamWinePath = filepath.Join("c:", "Program Files (x86)", "Steam")

func FindInstallationsWine(winePrefix string, launcher string, launchPath []string) ([]*common.Installation, []error) {
	platform := common.WineLauncherPlatform(winePrefix)

	if _, err := os.Stat(platform.ProcessPath(steamWinePath)); os.IsNotExist(err) {
		return nil, []error{fmt.Errorf("Steam is not installed in %s", winePrefix)}
	}

	return FindInstallationsSteam(
		steamWinePath,
		launcher,
		common.MakeLauncherPlatform(platform, func(_ string) []string { return launchPath }),
	)
}
