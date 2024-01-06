package launchers

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/epic"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/heroic"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/legendary"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/steam"
)

func GetInstallFinders() []common.InstallFinderFunc {
	return []common.InstallFinderFunc{
		epic.FindInstallations,
		heroic.FindInstallations,
		legendary.FindInstallations,
		steam.FindInstallations,
	}
}
