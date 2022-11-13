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

	heroicInstalls, heroicFindErrors := FindInstallationsWindowsHeroic()
	for _, install := range heroicInstalls {
		// Check if the install is already an Epic install
		found := false
		for _, existingInstall := range installs {
			if existingInstall.Path == install.Path {
				found = true
				break
			}
		}
		if !found {
			installs = append(installs, install)
		}
	}
	findErrors = append(findErrors, heroicFindErrors...)

	legendaryInstalls, legendaryFindErrors := FindInstallationsWindowsLegendary()
	for _, install := range legendaryInstalls {
		// Check if the install is already an Epic or Heroic install
		found := false
		for _, existingInstall := range installs {
			if existingInstall.Path == install.Path {
				found = true
				break
			}
		}
		if !found {
			installs = append(installs, install)
		}
	}
	findErrors = append(findErrors, legendaryFindErrors...)

	return installs, findErrors
}
