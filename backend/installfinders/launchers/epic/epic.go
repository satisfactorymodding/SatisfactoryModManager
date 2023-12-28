package epic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

type GameVersionFile struct {
	MajorVersion         int    `json:"MajorVersion"`
	MinorVersion         int    `json:"MinorVersion"`
	PatchVersion         int    `json:"PatchVersion"`
	Changelist           int    `json:"Changelist"`
	CompatibleChangelist int    `json:"CompatibleChangelist"`
	IsLicenseeVersion    int    `json:"IsLicenseeVersion"`
	IsPromotedBuild      int    `json:"IsPromotedBuild"`
	BranchName           string `json:"BranchName"`
	BuildID              string `json:"BuildId"`
}

type Manifest struct {
	CatalogNamespace string `json:"CatalogNamespace"`
	CatalogItemID    string `json:"CatalogItemID"`
	ManifestLocation string `json:"ManifestLocation"`
	InstallationGUID string `json:"InstallationGUID"`
	MainGameAppName  string `json:"MainGameAppName"`
	AppVersionString string `json:"AppVersionString"`
	InstallLocation  string `json:"InstallLocation"`
}

type GameManifest struct {
	AppName          string `json:"AppName"`
	CatalogNamespace string `json:"CatalogNamespace"`
	CatalogItemID    string `json:"CatalogItemID"`
}

var (
	EarlyAccessAppName  = "CrabEA"
	ExperimentalAppName = "CrabTest"
)

func GetEpicBranch(appName string) (common.GameBranch, error) {
	switch appName {
	case EarlyAccessAppName:
		return common.BranchEarlyAccess, nil
	case ExperimentalAppName:
		return common.BranchExperimental, nil
	default:
		return "", errors.New("unknown branch for " + appName)
	}
}

func findInstallationsEpic(epicManifestsPath string, launcher string, launchPath func(appName string) []string, processPath func(path string) string) ([]*common.Installation, []error) {
	if launchPath == nil {
		launchPath = func(appName string) []string { return nil }
	}

	if processPath == nil {
		processPath = func(path string) string { return path }
	}

	if _, err := os.Stat(epicManifestsPath); os.IsNotExist(err) {
		return nil, []error{errors.New("Epic is not installed")}
	}

	manifests, err := os.ReadDir(epicManifestsPath)
	if err != nil {
		return nil, []error{errors.Wrap(err, "Failed to list Epic manifests")}
	}

	installs := make([]*common.Installation, 0)
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

		var epicManifest Manifest
		if err := json.Unmarshal(manifestData, &epicManifest); err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to parse Epic manifest %s", manifestName))
			continue
		}

		if epicManifest.CatalogNamespace != "crab" {
			continue
		}

		installLocation := processPath(epicManifest.InstallLocation)

		gameManifestName := fmt.Sprintf("%s.mancpn", epicManifest.InstallationGUID)
		gameManifestPath := processPath(filepath.Join(epicManifest.ManifestLocation, gameManifestName))
		gameManifestData, err := os.ReadFile(gameManifestPath)
		if err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to read Epic game manifest %s", gameManifestName))
			continue
		}

		var epicGameManifest GameManifest
		if err := json.Unmarshal(gameManifestData, &epicGameManifest); err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to parse Epic game manifest %s", gameManifestName))
			continue
		}

		if epicGameManifest.CatalogNamespace != epicManifest.CatalogNamespace ||
			epicGameManifest.CatalogItemID != epicManifest.CatalogItemID ||
			epicGameManifest.AppName != epicManifest.MainGameAppName {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
				Inner: errors.New("Mismatching manifest data"),
			})
			continue
		}

		existingIdx := -1
		for i := range installs {
			if installs[i].Path == installLocation {
				existingIdx = i
				break
			}
		}

		if existingIdx != -1 {
			continue
		}

		versionFilePath := filepath.Join(installLocation, "Engine", "Binaries", "Win64", "FactoryGame-Win64-Shipping.version")
		if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
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

		branch, err := GetEpicBranch(epicManifest.MainGameAppName)
		if err != nil {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
				Inner: err,
			})
			continue
		}

		installs = append(installs, &common.Installation{
			Path:       filepath.Clean(installLocation),
			Version:    versionData.Changelist,
			Type:       common.InstallTypeWindowsClient,
			Location:   common.LocationTypeLocal,
			Branch:     branch,
			Launcher:   launcher,
			LaunchPath: launchPath(epicManifest.MainGameAppName),
		})
	}

	return installs, findErrors
}
