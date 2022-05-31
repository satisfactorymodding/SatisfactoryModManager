package windows

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/installFinders/types"
)

var EpicManifestsFolder = path.Join(os.Getenv("PROGRAMDATA"), "Epic", "EpicGamesLauncher", "Data", "Manifests")
var UEInstalledManifest = path.Join(os.Getenv("PROGRAMDATA"), "Epic", "UnrealEngineLauncher", "LauncherInstalled.dat")

func FindInstallationsEpic() ([]*types.Installation, []string, []error) {
	if _, err := os.Stat(EpicManifestsFolder); os.IsNotExist(err) {
		return nil, nil, []error{errors.New("Epic is not installed")}
	}

	manifests, err := ioutil.ReadDir(EpicManifestsFolder)
	if err != nil {
		return nil, nil, []error{errors.Wrap(err, "Failed to list Epic manifests")}
	}

	installs := []*types.Installation{}
	invalidInstalls := []string{}
	findErrors := []error{}

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

		var epicManifest types.EpicManifest
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

		var epicGameManifest types.EpicGameManifest
		if err := json.Unmarshal(gameManifestData, &epicGameManifest); err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to parse Epic game manifest %s", gameManifestName))
			continue
		}

		if epicGameManifest.CatalogNamespace != epicManifest.CatalogNamespace ||
			epicGameManifest.CatalogItemId != epicManifest.CatalogItemId ||
			epicGameManifest.AppName != epicManifest.MainGameAppName {
			invalidInstalls = append(invalidInstalls, epicManifest.InstallLocation)
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
			invalidInstalls = append(invalidInstalls, epicManifest.InstallLocation)
			continue
		}

		versionFile, err := os.ReadFile(versionFilePath)
		if err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to read version file %s", versionFilePath))
			continue
		}

		var versionData types.GameVersionFile
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
			invalidInstalls = append(invalidInstalls, epicManifest.InstallLocation)
			continue
		}

		installs = append(installs, &types.Installation{
			Path:       epicManifest.InstallLocation,
			Version:    versionData.Changelist,
			Branch:     branch,
			Launcher:   "Epic Games",
			LaunchPath: fmt.Sprintf(`start "" "com.epicgames.launcher://apps/%s?action=launch&silent=true"`, epicManifest.MainGameAppName),
		})
	}

	return installs, invalidInstalls, findErrors
}
