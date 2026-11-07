package epic

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("EpicGames", func() ([]*common.Installation, []error) {
		var errors []error

		var v1Path string
		programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, 0)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to get ProgramData folder for Epic V1: %w", err))
		} else {
			v1Path = filepath.Join(programData, "Epic")
		}

		var v2Path string
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to get user cache dir for Epic V2: %w", err))
		} else {
			v2Path = filepath.Join(cacheDir, "com.epicgames.eglv2")
		}

		installations, findErrors := FindInstallationsEpic(
			v1Path,
			v2Path,
			"Epic Games",
			common.MakeLauncherPlatform(
				common.NativePlatform(),
				func(appName string) []string {
					return []string{
						"cmd",
						"/C",
						`start`,
						``,
						// The extra space at the end is required for exec to escape the argument with double quotes
						// Otherwise, the & is interpreted as a command sequence
						`com.epicgames.launcher://apps/` + appName + `?action=launch&silent=true `,
					}
				},
			),
		)
		return installations, append(errors, findErrors...)
	})
}
