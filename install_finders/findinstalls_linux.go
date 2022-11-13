//go:build !bindings

package install_finders

func FindInstallations() ([]*Installation, []error) {
	installs := []*Installation{}
	findErrors := []error{}

	lutrisInstalls, lutrisFindErrors := FindInstallationsLinuxLutris()
	installs = append(installs, lutrisInstalls...)
	findErrors = append(findErrors, lutrisFindErrors...)

	lutrisFlatpakInstalls, lutrisFlatpakFindErrors := FindInstallationsLinuxLutrisFlatpak()
	installs = append(installs, lutrisFlatpakInstalls...)
	findErrors = append(findErrors, lutrisFlatpakFindErrors...)

	steaminstalls, steamfinderrors := FindInstallationsLinuxSteam()
	installs = append(installs, steaminstalls...)
	findErrors = append(findErrors, steamfinderrors...)

	steamFlatpakInstalls, steamFlatpakFindErrors := FindInstallationsLinuxSteamFlatpak()
	installs = append(installs, steamFlatpakInstalls...)
	findErrors = append(findErrors, steamFlatpakFindErrors...)

	heroicInstalls, heroicFindErrors := FindInstallationsLinuxHeroic()
	installs = append(installs, heroicInstalls...)
	findErrors = append(findErrors, heroicFindErrors...)

	heroicFlatpakInstalls, heroicFlatpakFindErrors := FindInstallationsLinuxHeroicFlatpak()
	installs = append(installs, heroicFlatpakInstalls...)
	findErrors = append(findErrors, heroicFlatpakFindErrors...)

	legendaryInstalls, legendaryFindErrors := FindInstallationsLinuxLegendary()
	for _, install := range legendaryInstalls {
		// Check if the install is already a Heroic install
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
