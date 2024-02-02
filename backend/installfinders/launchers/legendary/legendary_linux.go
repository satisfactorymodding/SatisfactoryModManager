package legendary

import (
	"fmt"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func FindInstallations() ([]*common.Installation, []error) {
	legendaryDataPath, err := getGlobalLegendaryDataPath("")
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get legendary config path: %w", err)}
	}
	return FindInstallationsIn(legendaryDataPath, "Legendary")
}
