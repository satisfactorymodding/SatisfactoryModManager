package epic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

type Manifest struct {
	CatalogNamespace string `json:"CatalogNamespace"`
	CatalogItemID    string `json:"CatalogItemID"`
	MainGameAppName  string `json:"MainGameAppName"`
	AppName          string `json:"AppName"`
	AppVersionString string `json:"AppVersionString"`
	InstallLocation  string `json:"InstallLocation"`
}

type EOSInstall struct {
	V4 struct {
		State      string `json:"state"`
		Dir        string `json:"dir"`
		SandboxID  string `json:"sandboxId"`
		ArtifactID string `json:"artifactId"`
	} `json:"v4"`
}

type V2Installation struct {
	InstallationID string `json:"installationId"`
	Complete       bool   `json:"complete"`
	Deleted        bool   `json:"deleted,omitempty"`
	Dir            string `json:"dir"`
	Namespace      string `json:"namespace"`
	ArtifactID     string `json:"artifactId"`
}

var (
	EarlyAccessAppName                 = "CrabEA"
	ExperimentalAppName                = "CrabTest"
	EarlyAccessDedicatedServerAppName  = "CrabDedicatedServer"
	ExperimentalDedicatedServerAppName = "c509233193024c5f8124467d3aa36199"
	CatalogNamespace                   = "crab"
)

func GetEpicBranch(appName string) (common.GameBranch, error) {
	switch appName {
	case EarlyAccessAppName:
		return common.BranchStable, nil
	case ExperimentalAppName:
		return common.BranchExperimental, nil
	case EarlyAccessDedicatedServerAppName:
		return common.BranchStable, nil
	case ExperimentalDedicatedServerAppName:
		return common.BranchExperimental, nil
	default:
		return "", fmt.Errorf("unknown branch for '%s'", appName)
	}
}

func FindInstallationsEpic(epicV1DataPath string, epicV2DataPath string, launcher string, platform common.LauncherPlatform) ([]*common.Installation, []error) {
	installs := make([]*common.Installation, 0)
	var findErrors []error

	if epicV1DataPath != "" {
		v1Installs, v1Errors := FindInstallationsEpicV1(epicV1DataPath, launcher, platform)
		installs = append(installs, v1Installs...)
		findErrors = append(findErrors, v1Errors...)
	}

	if epicV2DataPath != "" {
		v2Installs, v2Errors := FindInstallationsEpicV2(epicV2DataPath, launcher, platform)
		installs = append(installs, v2Installs...)
		findErrors = append(findErrors, v2Errors...)
	}

	return installs, findErrors
}

var (
	errEpicV1ManifestsNotExist = fmt.Errorf("epic v1 manifests do not exist")
	errEpicV1EOSNotExist       = fmt.Errorf("epic v1 install states do not exist")
)

func FindInstallationsEpicV1(epicV1DataPath string, launcher string, platform common.LauncherPlatform) ([]*common.Installation, []error) {
	manifestsPath := filepath.Join(epicV1DataPath, "EpicGamesLauncher", "Data", "Manifests")
	eosPath := filepath.Join(epicV1DataPath, "EpicOnlineServicesShared", "InstallHelper", "InstalledItems")

	installsManifests, errorsManifests := FindInstallationsEpicV1Manifest(manifestsPath, launcher, platform)
	installsEOS, errorsEOS := FindInstallationsEpicV1EOS(eosPath, launcher, platform)

	noManifests := false
	noEOS := false
	if len(installsManifests) == 0 && len(errorsManifests) == 1 && errors.Is(errorsManifests[0], errEpicV1ManifestsNotExist) {
		noManifests = true
		errorsManifests = []error{}
	}
	if len(installsEOS) == 0 && len(errorsEOS) == 1 && errors.Is(errorsEOS[0], errEpicV1EOSNotExist) {
		noEOS = true
		errorsEOS = []error{}
	}

	if noManifests && noEOS {
		return nil, []error{fmt.Errorf("epic v1 is not installed")}
	}

	return append(installsManifests, installsEOS...), append(errorsManifests, errorsEOS...)
}

func FindInstallationsEpicV1Manifest(epicManifestsPath string, launcher string, platform common.LauncherPlatform) ([]*common.Installation, []error) {
	manifestsDir := platform.ProcessPath(epicManifestsPath)
	if _, err := os.Stat(manifestsDir); os.IsNotExist(err) {
		return nil, []error{errEpicV1ManifestsNotExist}
	}

	manifests, err := os.ReadDir(manifestsDir)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to list Epic manifests: %w", err)}
	}

	installs := make([]*common.Installation, 0)
	var findErrors []error

	for _, manifest := range manifests {
		manifestName := manifest.Name()
		manifestPath := platform.ProcessPath(filepath.Join(epicManifestsPath, manifestName))

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

		if epicManifest.CatalogNamespace != CatalogNamespace {
			continue
		}

		installLocation := platform.ProcessPath(epicManifest.InstallLocation)

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

		// Epic V1 can only launch games of the same platform
		gamePlatform := platform.Platform

		installType, version, savedPath, err := common.GetGameInfo(installLocation, gamePlatform)
		if err != nil {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
				Inner: err,
			})
			continue
		}

		appName := epicManifest.MainGameAppName
		if appName == "" {
			appName = epicManifest.AppName
		}

		branch, err := GetEpicBranch(appName)
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
			LaunchPath: platform.LauncherCommand(appName),
			SavedPath:  savedPath,
		})
	}

	return installs, findErrors
}

func FindInstallationsEpicV1EOS(eosInstallStatePath string, launcher string, platform common.LauncherPlatform) ([]*common.Installation, []error) {
	installStatesDir := platform.ProcessPath(eosInstallStatePath)
	if _, err := os.Stat(installStatesDir); os.IsNotExist(err) {
		return nil, []error{errEpicV1EOSNotExist}
	}

	installStates, err := os.ReadDir(installStatesDir)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to list Epic install states: %w", err)}
	}

	installs := make([]*common.Installation, 0)
	var findErrors []error

	for _, installState := range installStates {
		installStateName := installState.Name()
		installStatePath := platform.ProcessPath(filepath.Join(eosInstallStatePath, installStateName))

		if fileInfo, err := os.Stat(installStatePath); os.IsNotExist(err) || fileInfo.IsDir() {
			continue
		}

		if !strings.HasSuffix(installStateName, ".egi") {
			continue
		}

		installStateData, err := os.ReadFile(installStatePath)
		if err != nil {
			findErrors = append(findErrors, fmt.Errorf("failed to read Epic install state %s: %w", installStateName, err))
			continue
		}

		var eosInstall EOSInstall
		if err := json.Unmarshal(installStateData, &eosInstall); err != nil {
			findErrors = append(findErrors, fmt.Errorf("failed to parse Epic install state %s: %w", installStateName, err))
			continue
		}

		if eosInstall.V4.SandboxID != CatalogNamespace {
			continue
		}

		installLocation := platform.ProcessPath(eosInstall.V4.Dir)

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

		// Epic V1 can only launch games of the same platform
		gamePlatform := platform.Platform

		installType, version, savedPath, err := common.GetGameInfo(installLocation, gamePlatform)
		if err != nil {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
				Inner: err,
			})
			continue
		}

		appName := eosInstall.V4.ArtifactID

		branch, err := GetEpicBranch(appName)
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
			LaunchPath: platform.LauncherCommand(appName),
			SavedPath:  savedPath,
		})
	}

	return installs, findErrors
}

func FindInstallationsEpicV2(epicV2DataPath string, launcher string, platform common.LauncherPlatform) ([]*common.Installation, []error) {
	installationsFile := filepath.Join(epicV2DataPath, "installations.jsonl")

	if _, err := os.Stat(installationsFile); os.IsNotExist(err) {
		return nil, []error{fmt.Errorf("epic v2 is not installed")}
	}

	f, err := os.Open(installationsFile)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to open Epic v2 installations file: %w", err)}
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to read Epic v2 installations file: %w", err)}
	}

	lines := strings.Split(string(data), "\n")

	installsMap := make(map[string]V2Installation)

	for _, line := range lines {
		var entry V2Installation
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		if entry.InstallationID == "" {
			continue
		}

		if entry.Deleted {
			delete(installsMap, entry.InstallationID)
			continue
		}

		installsMap[entry.InstallationID] = entry
	}

	installs := make([]*common.Installation, 0)
	var findErrors []error

	for _, entry := range installsMap {
		if entry.Namespace != CatalogNamespace {
			continue
		}

		installLocation := platform.ProcessPath(entry.Dir)

		existingIdx := -1
		for i := range installs {
			if common.OsPathEqual(installs[i].Path, installLocation) {
				existingIdx = i
				break
			}
		}

		if existingIdx != -1 {
			continue
		}

		gamePlatform := platform.Platform
		if platform.Os() != "windows" {
			gamePlatform = common.WineLauncherPlatform(filepath.Join(epicV2DataPath, "compat", "compatdata", entry.ArtifactID, "pfx"))
		}

		installType, version, savedPath, err := common.GetGameInfo(installLocation, gamePlatform)
		if err != nil {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
				Inner: err,
			})
			continue
		}

		appName := entry.ArtifactID

		branch, err := GetEpicBranch(appName)
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
			LaunchPath: platform.LauncherCommand(appName),
			SavedPath:  savedPath,
		})
	}

	return installs, findErrors
}
