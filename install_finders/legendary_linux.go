package install_finders

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

func FindInstallationsLinuxLegendary() ([]*Installation, []error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get user home dir")}
	}
	legendaryInstalledPath := filepath.Join(homeDir, ".config", "legendary", "installed.json")
	if _, err := os.Stat(legendaryInstalledPath); os.IsNotExist(err) {
		return nil, []error{errors.New("legendary not installed")}
	}
	var legendaryData LegendaryData
	legendaryDataFile, err := os.ReadFile(legendaryInstalledPath)
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to read legendary installed.json")}
	}
	err = json.Unmarshal(legendaryDataFile, &legendaryData)
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to parse legendary installed.json output")}
	}

	_, err = exec.LookPath("legendary")
	canLaunchLegendary := err == nil

	var installs []*Installation
	var findErrors []error

	for _, legendaryGame := range legendaryData {
		var branch GameBranch
		switch legendaryGame.AppName[4:] {
		case "EA":
			branch = BRANCH_EARLY_ACCESS
		case "Test":
			branch = BRANCH_EXPERIMENTAL
		default:
			findErrors = append(findErrors, errors.New("unknown branch for "+legendaryGame.AppName))
			continue
		}
		version, err := strconv.Atoi(legendaryGame.Version)
		if err != nil {
			findErrors = append(findErrors, errors.Wrap(err, "failed to parse version for "+legendaryGame.AppName))
			continue
		}
		var launchPath []string
		if canLaunchLegendary {
			launchPath = []string{"legendary", "launch", legendaryGame.AppName}
		}
		installs = append(installs, &Installation{
			Path:       legendaryGame.InstallPath,
			Version:    version,
			Branch:     branch,
			Launcher:   "Legendary",
			LaunchPath: launchPath,
		})
	}
	return installs, findErrors
}
