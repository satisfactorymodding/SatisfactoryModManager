package launchers

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/heroic"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/legendary"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/lutris"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/steam"
)

func GetInstallFinders() []common.InstallFinderFunc {
	return []common.InstallFinderFunc{
		heroic.FindInstallations,
		legendary.FindInstallations,
		lutris.FindInstallations,
		steam.FindInstallations,
	}
}
