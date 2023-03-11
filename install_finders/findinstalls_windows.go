//go:build !bindings

package install_finders

func FindInstallations() ([]*Installation, []error) {
	// Check Heroic before Legendary, since Heroic installs will show up in Legendary too
	// Check Epic before Heroic, since Heroic can use Epic installs
	return FindAll(FindInstallationsWindowsEpic, FindInstallationsWindowsSteam, FindInstallationsWindowsHeroic, FindInstallationsWindowsLegendary)
}
