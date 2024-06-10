package epic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

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
	EarlyAccessAppName                 = "CrabEA"
	ExperimentalAppName                = "CrabTest"
	EarlyAccessDedicatedServerAppName  = "CrabDedicatedServer"
	ExperimentalDedicatedServerAppName = "c509233193024c5f8124467d3aa36199"
)

func GetEpicBranch(appName string) (common.GameBranch, error) {
	switch appName {
	case EarlyAccessAppName:
		return common.BranchEarlyAccess, nil
	case ExperimentalAppName:
		return common.BranchExperimental, nil
	case EarlyAccessDedicatedServerAppName:
		return common.BranchEarlyAccess, nil
	case ExperimentalDedicatedServerAppName:
		return common.BranchExperimental, nil
	default:
		return "", fmt.Errorf("unknown branch for " + appName)
	}
}

func FindInstallationsEpic(epicManifestsPath string, launcher string, platform common.LauncherPlatform) ([]*common.Installation, []error) {
	if _, err := os.Stat(platform.ProcessPath(epicManifestsPath)); os.IsNotExist(err) {
		return nil, []error{fmt.Errorf("epic is not installed")}
	}

	manifests, err := os.ReadDir(platform.ProcessPath(epicManifestsPath))
	if err != nil {
		return nil, []error{fmt.Errorf("failed to list Epic manifests: %w", err)}
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
			findErrors = append(findErrors, fmt.Errorf("failed to read Epic manifest %s: %w", manifestName, err))
			continue
		}

		var epicManifest Manifest
		if err := json.Unmarshal(manifestData, &epicManifest); err != nil {
			findErrors = append(findErrors, fmt.Errorf("failed to parse Epic manifest %s: %w", manifestName, err))
			continue
		}

		if epicManifest.CatalogNamespace != "crab" {
			continue
		}

		installLocation := platform.ProcessPath(epicManifest.InstallLocation)

		gameManifestName := fmt.Sprintf("%s.mancpn", epicManifest.InstallationGUID)
		gameManifestPath := platform.ProcessPath(filepath.Join(epicManifest.ManifestLocation, gameManifestName))
		gameManifestData, err := os.ReadFile(gameManifestPath)
		if err != nil {
			findErrors = append(findErrors, fmt.Errorf("failed to read Epic game manifest %s: %w", gameManifestName, err))
			continue
		}

		var epicGameManifest GameManifest
		if err := json.Unmarshal(gameManifestData, &epicGameManifest); err != nil {
			findErrors = append(findErrors, fmt.Errorf("failed to parse Epic game manifest %s: %w", gameManifestName, err))
			continue
		}

		if epicGameManifest.CatalogNamespace != epicManifest.CatalogNamespace ||
			epicGameManifest.CatalogItemID != epicManifest.CatalogItemID ||
			epicGameManifest.AppName != epicManifest.MainGameAppName {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
				Inner: fmt.Errorf("mismatching manifest data"),
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

		// Epic can only launch games of the same platform
		gamePlatform := platform.Platform

		installType, version, savedPath, err := common.GetGameInfo(installLocation, gamePlatform)
		if err != nil {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
				Inner: err,
			})
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
			Version:    version,
			Type:       installType,
			Location:   common.LocationTypeLocal,
			Branch:     branch,
			Launcher:   launcher,
			LaunchPath: platform.LauncherCommand(epicManifest.MainGameAppName),
			SavedPath:  savedPath,
		})
	}

	return installs, findErrors
}
