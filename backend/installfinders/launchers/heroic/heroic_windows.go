package heroic

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func FindInstallations() ([]*common.Installation, []error) {
	return findInstallationsHeroic(false, "", "Heroic")
}
