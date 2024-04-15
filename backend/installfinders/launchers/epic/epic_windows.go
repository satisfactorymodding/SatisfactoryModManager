package epic

import (
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

var epicManifestsFolder = filepath.Join(os.Getenv("PROGRAMDATA"), "Epic", "EpicGamesLauncher", "Data", "Manifests")

func init() {
	launchers.Add("EpicGames", func() ([]*common.Installation, []error) {
		return FindInstallationsEpic(epicManifestsFolder, "Epic Games", func(appName string) []string {
			return []string{
				"cmd",
				"/C",
				`start`,
				``,
				// The extra space at the end is required for exec to escape the argument with double quotes
				// Otherwise, the & is interpreted as a command sequence
				`com.epicgames.launcher://apps/` + appName + `?action=launch&silent=true `,
			}
		}, nil)
	})
}
