package legendary

import (
	"fmt"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Legendary", func() ([]*common.Installation, []error) {
		legendaryDataPath, err := getGlobalLegendaryDataPath("")
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get legendary config path: %w", err)}
		}
		return FindInstallationsIn(legendaryDataPath, "Legendary")
	})
}
