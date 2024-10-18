package epic

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

var epicProgramDataManifestsFolder = filepath.Join("Epic", "EpicGamesLauncher", "Data", "Manifests")

func init() {
	launchers.Add("EpicGames", func() ([]*common.Installation, []error) {
		programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, 0)
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get ProgramData folder: %w", err)}
		}

		return FindInstallationsEpic(
			filepath.Join(programData, epicProgramDataManifestsFolder),
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
	})
}
