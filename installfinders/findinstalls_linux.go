package installfinders

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/installfinders/types"
)

func FindInstallations() ([]*types.Installation, []string, []error) {
	installs := []*types.Installation{}
	invalidInstalls := []string{}
	findErrors := []error{}

	return installs, invalidInstalls, findErrors
}
