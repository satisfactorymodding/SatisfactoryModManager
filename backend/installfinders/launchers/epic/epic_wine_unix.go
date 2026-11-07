//go:build unix

package epic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

var epicV1WineDataPath = filepath.Join("c:", "ProgramData", "Epic")

func FindInstallationsWine(winePrefix string, launcher string, launchPath []string) ([]*common.Installation, []error) {
	var errors []error
	platform := common.WineLauncherPlatform(winePrefix)

	var v2Path string
	cacheDir, err := platform.CacheDir()
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to get cache dir: %w", err))
	} else {
		v2Path = filepath.Join(cacheDir, "com.epicgames.eglv2")
	}

	hasV1 := true
	hasV2 := true
	if _, err := os.Stat(platform.ProcessPath(epicV1WineDataPath)); os.IsNotExist(err) {
		hasV1 = false
	}
	if _, err := os.Stat(platform.ProcessPath(v2Path)); os.IsNotExist(err) {
		hasV2 = false
	}

	if !hasV1 && !hasV2 {
		return nil, []error{fmt.Errorf("epic is not installed in %s", winePrefix)}
	}

	installs, findErrors := FindInstallationsEpic(
		epicV1WineDataPath,
		v2Path,
		launcher,
		common.MakeLauncherPlatform(platform, func(_ string) []string { return launchPath }),
	)

	return installs, append(errors, findErrors...)
}
