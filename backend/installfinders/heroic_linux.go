package installfinders

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

func FindInstallationsLinuxHeroic() ([]*Installation, []error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get user home dir")}
	}
	heroicConfigPath := filepath.Join(homeDir, ".config", "heroic", "config.json")
	heroicGamesConfigPath := filepath.Join(homeDir, ".config", "heroic", "GamesConfig")
	legendaryInstalledPath := filepath.Join(homeDir, ".config", "legendary", "installed.json")
	if _, err := os.Stat(heroicConfigPath); os.IsNotExist(err) {
		return nil, []error{errors.New("heroic not installed")}
	}
	if _, err := os.Stat(legendaryInstalledPath); os.IsNotExist(err) {
		return nil, []error{errors.New("heroic - legendary not installed")}
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

	installs := make([]*Installation, 0)
	var findErrors []error

	for key, legendaryGame := range legendaryData {
		if _, err := os.Stat(filepath.Join(heroicGamesConfigPath, key+".json")); os.IsNotExist(err) {
			// Not a heroic game
			continue
		}
		var branch GameBranch
		switch legendaryGame.AppName[4:] {
		case "EA":
			branch = BranchEarlyAccess
		case "Test":
			branch = BranchExperimental
		default:
			findErrors = append(findErrors, errors.New("unknown branch for "+legendaryGame.AppName))
			continue
		}
		version, err := strconv.Atoi(legendaryGame.Version)
		if err != nil {
			findErrors = append(findErrors, errors.Wrap(err, "failed to parse version for "+legendaryGame.AppName))
			continue
		}
		installs = append(installs, &Installation{
			Path:       legendaryGame.InstallPath,
			Version:    version,
			Type:       InstallTypeWindowsClient,
			Location:   LocationTypeLocal,
			Branch:     branch,
			Launcher:   "Heroic",
			LaunchPath: nil, // Heroic doesn't support launching from command line
		})
	}
	return installs, findErrors
}
