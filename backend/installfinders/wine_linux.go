package installfinders

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var epicManifestRelativePath = filepath.Join("Epic", "EpicGamesLauncher", "Data", "Manifests")

func findInstallationsWineEpic(winePrefix string, launcher string, launchPath []string) ([]*Installation, []error) {
	wineWindowsRoot := filepath.Join(winePrefix, "dosdevices")
	epicManifestsPath := filepath.Join(wineWindowsRoot, "c:", "ProgramData", epicManifestRelativePath)
	if _, err := os.Stat(epicManifestsPath); os.IsNotExist(err) {
		return nil, []error{errors.New("Epic is not installed in " + winePrefix)}
	}

	manifests, err := os.ReadDir(epicManifestsPath)
	if err != nil {
		return nil, []error{errors.Wrap(err, "Failed to list Epic manifests in "+winePrefix)}
	}

	installs := make([]*Installation, 0)
	var findErrors []error

	for _, manifest := range manifests {
		manifestName := manifest.Name()
		manifestPath := filepath.Join(epicManifestsPath, manifestName)

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

		linuxInstallLocation := strings.ToLower(epicManifest.InstallLocation[0:1]) + strings.ReplaceAll(epicManifest.InstallLocation[1:], "\\", "/")
		wineInstallLocation := filepath.Join(wineWindowsRoot, linuxInstallLocation)

		gameManifestName := fmt.Sprintf("%s.mancpn", epicManifest.InstallationGUID)
		gameManifestPath := filepath.Join(epicManifest.ManifestLocation, gameManifestName)
		linuxGameManifestPath := strings.ToLower(gameManifestPath[0:1]) + strings.ReplaceAll(gameManifestPath[1:], "\\", "/")
		wineGameManifestPath := filepath.Join(wineWindowsRoot, linuxGameManifestPath)
		gameManifestData, err := os.ReadFile(wineGameManifestPath)
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
			epicGameManifest.CatalogItemID != epicManifest.CatalogItemID ||
			epicGameManifest.AppName != epicManifest.MainGameAppName {
			findErrors = append(findErrors, InstallFindError{
				Path:  wineInstallLocation,
				Inner: errors.New("Mismatching manifest data"),
			})
			continue
		}

		existingIdx := -1
		for i := range installs {
			if installs[i].Path == wineInstallLocation {
				existingIdx = i
				break
			}
		}

		if existingIdx != -1 {
			continue
		}

		versionFilePath := filepath.Join(wineInstallLocation, "Engine", "Binaries", "Win64", "FactoryGame-Win64-Shipping.version")
		if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
			findErrors = append(findErrors, InstallFindError{
				Path:  wineInstallLocation,
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

		var branch GameBranch
		if epicManifest.MainGameAppName == "CrabEA" {
			branch = BranchEarlyAccess
		} else if epicManifest.MainGameAppName == "CrabTest" {
			branch = BranchExperimental
		} else {
			findErrors = append(findErrors, InstallFindError{
				Path:  wineInstallLocation,
				Inner: errors.New("Invalid branch " + epicManifest.MainGameAppName),
			})
			continue
		}

		installs = append(installs, &Installation{
			Path:       filepath.Clean(wineInstallLocation),
			Version:    versionData.Changelist,
			Type:       InstallTypeWindowsClient,
			Location:   LocationTypeLocal,
			Branch:     branch,
			Launcher:   launcher,
			LaunchPath: launchPath,
		})
	}

	return installs, findErrors
}
