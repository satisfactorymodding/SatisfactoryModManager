package installFinders

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/installFinders/types"
	"github.com/satisfactorymodding/SatisfactoryModManager/installFinders/windows"
)

func FindInstallations() ([]*types.Installation, []string, []error) {
	installs := []*types.Installation{}
	invalidInstalls := []string{}
	findErrors := []error{}

	epicInstalls, epicInvalidInstalls, epicFindErrors := windows.FindInstallationsEpic()
	installs = append(installs, epicInstalls...)
	invalidInstalls = append(invalidInstalls, epicInvalidInstalls...)
	findErrors = append(findErrors, epicFindErrors...)

	steamInstalls, steamInvalidInstalls, steamFindErrors := windows.FindInstallationsSteam()
	installs = append(installs, steamInstalls...)
	invalidInstalls = append(invalidInstalls, steamInvalidInstalls...)
	findErrors = append(findErrors, steamFindErrors...)

	return installs, invalidInstalls, findErrors
}
