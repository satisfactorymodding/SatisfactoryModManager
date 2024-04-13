package launchers

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/heroic"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/legendary"
)

func GetInstallFinders() []common.InstallFinderFunc {
	return []common.InstallFinderFunc{
		heroic.FindInstallations,
		legendary.FindInstallations,
	}
}
