package installfinders

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/installfinders/types"
	"github.com/satisfactorymodding/SatisfactoryModManager/installfinders/wininstalls"
)

func FindInstallations() ([]*types.Installation, []string, []error) {
	installs := []*types.Installation{}
	invalidInstalls := []string{}
	findErrors := []error{}

	epicInstalls, epicInvalidInstalls, epicFindErrors := wininstalls.FindInstallationsEpic()
	installs = append(installs, epicInstalls...)
	invalidInstalls = append(invalidInstalls, epicInvalidInstalls...)
	findErrors = append(findErrors, epicFindErrors...)

	steamInstalls, steamInvalidInstalls, steamFindErrors := wininstalls.FindInstallationsSteam()
	installs = append(installs, steamInstalls...)
	invalidInstalls = append(invalidInstalls, steamInvalidInstalls...)
	findErrors = append(findErrors, steamFindErrors...)

	return installs, invalidInstalls, findErrors
}
