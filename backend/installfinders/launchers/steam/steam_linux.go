package steam

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func FindInstallations() ([]*common.Installation, []error) {
	return common.FindAll(findInstallationsNative, findInstallationsFlatpak)
}

func findInstallationsNative() ([]*common.Installation, []error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get user home dir")}
	}
	steamPath := filepath.Join(homeDir, ".steam", "steam")
	if _, err := os.Stat(steamPath); os.IsNotExist(err) {
		return nil, []error{errors.New("steam not installed")}
	}
	return findInstallationsSteam(
		steamPath,
		"Steam",
		[]string{
			"steam",
		},
	)
}

func findInstallationsFlatpak() ([]*common.Installation, []error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get user home dir")}
	}

	steamPath := filepath.Join(homeDir, ".var", "app", "com.valvesoftware.Steam", ".steam", "steam")
	if _, err := os.Stat(steamPath); os.IsNotExist(err) {
		return nil, []error{errors.New("steam-flatpak not installed")}
	}
	return findInstallationsSteam(
		steamPath,
		"Steam",
		[]string{
			"flatpak",
			"run",
			"com.valvesoftware.Steam",
		},
	)
}
