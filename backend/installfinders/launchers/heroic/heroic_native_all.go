package heroic

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Heroic", func() ([]*common.Installation, []error) {
		return findInstallationsHeroic(false, "", "Heroic")
	})
}
