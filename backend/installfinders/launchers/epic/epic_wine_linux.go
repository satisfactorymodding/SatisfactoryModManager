package epic

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

var epicWineManifestPath = filepath.Join("c:", "ProgramData", "Epic", "EpicGamesLauncher", "Data", "Manifests")

func FindInstallationsWine(winePrefix string, launcher string, launchPath []string) ([]*common.Installation, []error) {
	wineWindowsRoot := filepath.Join(winePrefix, "dosdevices")
	epicManifestsPath := filepath.Join(wineWindowsRoot, epicWineManifestPath)

	if _, err := os.Stat(epicManifestsPath); os.IsNotExist(err) {
		return nil, []error{errors.New("Epic is not installed in " + winePrefix)}
	}

	return findInstallationsEpic(epicManifestsPath, launcher, func(appName string) []string { return launchPath }, func(path string) string {
		return filepath.Join(wineWindowsRoot, strings.ToLower(path[0:1])+strings.ReplaceAll(path[1:], "\\", "/"))
	})
}
