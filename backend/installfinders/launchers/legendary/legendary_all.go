package legendary

import (
	"fmt"
	"os/exec"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Legendary", func() ([]*common.Installation, []error) {
		legendaryDataPath, err := getGlobalLegendaryDataPath("")
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get legendary config path: %w", err)}
		}

		_, err = exec.LookPath("legendary")
		canLaunchLegendary := err == nil

		return FindInstallationsIn(
			legendaryDataPath,
			"Legendary",
			nil,
			common.MakeLauncherPlatform(
				common.NativePlatform(),
				func(appName string) []string {
					if !canLaunchLegendary {
						return nil
					}
					return []string{"legendary", "launch", appName}
				},
			),
		)
	})
}
