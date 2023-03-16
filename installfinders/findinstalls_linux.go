//go:build !bindings

package installfinders

func FindInstallations() ([]*Installation, []error) {
	// Check Heroic before Legendary, since Heroic installs will show up in Legendary too
	return FindAll(FindInstallationsLinuxLutris, FindInstallationsLinuxLutrisFlatpak, FindInstallationsLinuxSteam, FindInstallationsLinuxSteamFlatpak, FindInstallationsLinuxHeroic, FindInstallationsLinuxHeroicFlatpak, FindInstallationsLinuxLegendary)
}
