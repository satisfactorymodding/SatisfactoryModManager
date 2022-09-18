package install_finders

func FindInstallations() ([]*Installation, []error) {
	installs := []*Installation{}
	findErrors := []error{}

	epicInstalls, epicFindErrors := FindInstallationsWindowsEpic()
	installs = append(installs, epicInstalls...)
	findErrors = append(findErrors, epicFindErrors...)

	steamInstalls, steamFindErrors := FindInstallationsWindowsSteam()
	installs = append(installs, steamInstalls...)
	findErrors = append(findErrors, steamFindErrors...)

	return installs, findErrors
}
