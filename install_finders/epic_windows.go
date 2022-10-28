package install_finders

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
)

var EpicManifestsFolder = path.Join(os.Getenv("PROGRAMDATA"), "Epic", "EpicGamesLauncher", "Data", "Manifests")

func FindInstallationsWindowsEpic() ([]*Installation, []error) {
	if _, err := os.Stat(EpicManifestsFolder); os.IsNotExist(err) {
		return nil, []error{errors.New("Epic is not installed")}
	}

	manifests, err := os.ReadDir(EpicManifestsFolder)
	if err != nil {
		return nil, []error{errors.Wrap(err, "Failed to list Epic manifests")}
	}

	var installs []*Installation
	var findErrors []error

	for _, manifest := range manifests {
		manifestName := manifest.Name()
		manifestPath := path.Join(EpicManifestsFolder, manifestName)

		if fileInfo, err := os.Stat(manifestPath); os.IsNotExist(err) || fileInfo.IsDir() {
			continue
		}

		manifestData, err := os.ReadFile(manifestPath)
		if err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to read Epic manifest %s", manifestName))
			continue
		}

		var epicManifest EpicManifest
		if err := json.Unmarshal(manifestData, &epicManifest); err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to parse Epic manifest %s", manifestName))
			continue
		}

		if epicManifest.CatalogNamespace != "crab" {
			continue
		}

		gameManifestName := fmt.Sprintf("%s.mancpn", epicManifest.InstallationGuid)
		gameManifestPath := path.Join(epicManifest.ManifestLocation, gameManifestName)
		gameManifestData, err := os.ReadFile(gameManifestPath)
		if err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to read Epic game manifest %s", gameManifestName))
			continue
		}

		var epicGameManifest EpicGameManifest
		if err := json.Unmarshal(gameManifestData, &epicGameManifest); err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to parse Epic game manifest %s", gameManifestName))
			continue
		}

		if epicGameManifest.CatalogNamespace != epicManifest.CatalogNamespace ||
			epicGameManifest.CatalogItemId != epicManifest.CatalogItemId ||
			epicGameManifest.AppName != epicManifest.MainGameAppName {
			findErrors = append(findErrors, InstallFindError{
				Path:  epicManifest.InstallLocation,
				Inner: errors.New("Mismatching manifest data"),
			})
			continue
		}

		existingIdx := -1
		for i := range installs {
			if installs[i].Path == epicManifest.InstallLocation {
				existingIdx = i
				break
			}
		}

		if existingIdx != -1 {
			continue
		}

		versionFilePath := path.Join(epicManifest.InstallLocation, "Engine", "Binaries", "Win64", "FactoryGame-Win64-Shipping.version")
		if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
			findErrors = append(findErrors, InstallFindError{
				Path:  epicManifest.InstallLocation,
				Inner: errors.Wrap(err, "failed to read game version"),
			})
			continue
		}

		versionFile, err := os.ReadFile(versionFilePath)
		if err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to read version file %s", versionFilePath))
			continue
		}

		var versionData GameVersionFile
		if err := json.Unmarshal(versionFile, &versionData); err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to parse version file %s", versionFilePath))
			continue
		}

		var branch string
		if epicManifest.MainGameAppName == "CrabEA" {
			branch = "Early Access"
		} else if epicManifest.MainGameAppName == "CrabTest" {
			branch = "Experimental"
		} else {
			findErrors = append(findErrors, InstallFindError{
				Path:  epicManifest.InstallLocation,
				Inner: errors.New("Invalid branch " + epicManifest.MainGameAppName),
			})
			continue
		}

		installs = append(installs, &Installation{
			Path:     epicManifest.InstallLocation,
			Version:  versionData.Changelist,
			Branch:   branch,
			Launcher: "Epic Games",
			LaunchPath: []string{
				"cmd",
				"/C",
				`start`,
				``,
				// The extra space at the end is required for exec to escape the argument with double quotes
				// Otherwise, the & is interpreted as a command sequence
				`com.epicgames.launcher://apps/` + epicManifest.MainGameAppName + `?action=launch&silent=true `,
			},
		})
	}

	return installs, findErrors
}
