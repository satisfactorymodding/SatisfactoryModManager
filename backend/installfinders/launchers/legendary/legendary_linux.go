package legendary

import (
	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func FindInstallations() ([]*common.Installation, []error) {
	legendaryDataPath, err := getGlobalLegendaryDataPath("")
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get legendary config path")}
	}
	return FindInstallationsIn(legendaryDataPath, "Legendary")
}
