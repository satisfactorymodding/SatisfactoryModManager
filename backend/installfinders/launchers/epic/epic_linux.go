package epic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Epic Games", func() ([]*common.Installation, []error) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
		}

		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			dataHome = filepath.Join(homeDir, ".local", "share")
		}

		return FindInstallationsEpicV2(
			filepath.Join(dataHome, "com.epicgames.eglv2"),
			"Epic Games",
			common.MakeLauncherPlatform(
				common.NativePlatform(),
				func(appName string) []string {
					return []string{
						"xdg-open",
						"com.epicgames.launcher://apps/" + appName + "?action=launch&silent=true",
					}
				},
			),
		)
	})
}
