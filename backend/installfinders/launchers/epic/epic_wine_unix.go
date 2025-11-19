//go:build unix

package epic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

var epicWineManifestPath = filepath.Join("c:", "ProgramData", "Epic", "EpicGamesLauncher", "Data", "Manifests")

func FindInstallationsWine(winePrefix string, launcher string, launchPath []string) ([]*common.Installation, []error) {
	platform := common.WineLauncherPlatform(winePrefix)

	if _, err := os.Stat(platform.ProcessPath(epicWineManifestPath)); os.IsNotExist(err) {
		return nil, []error{fmt.Errorf("epic is not installed in %s", winePrefix)}
	}

	return FindInstallationsEpic(
		epicWineManifestPath,
		launcher,
		common.MakeLauncherPlatform(platform, func(_ string) []string { return launchPath }),
	)
}
