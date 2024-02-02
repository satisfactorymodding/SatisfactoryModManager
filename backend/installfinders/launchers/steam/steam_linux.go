package steam

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func FindInstallations() ([]*common.Installation, []error) {
	return common.FindAll(findInstallationsNative, findInstallationsFlatpak)
}

func findInstallationsNative() ([]*common.Installation, []error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
	}
	steamPath := filepath.Join(homeDir, ".steam", "steam")
	if _, err := os.Stat(steamPath); os.IsNotExist(err) {
		return nil, []error{fmt.Errorf("steam not installed")}
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
		return nil, []error{fmt.Errorf("failed to get user home dir: %w", err)}
	}

	steamPath := filepath.Join(homeDir, ".var", "app", "com.valvesoftware.Steam", ".steam", "steam")
	if _, err := os.Stat(steamPath); os.IsNotExist(err) {
		return nil, []error{fmt.Errorf("steam-flatpak not installed")}
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
