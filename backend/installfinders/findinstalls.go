//go:build !bindings

package installfinders

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func FindInstallations() ([]*common.Installation, []error) {
	return common.FindAll(launchers.GetInstallFinders()...)
}
