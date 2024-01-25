package legendary

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/epic"
)

type Game struct {
	AppName           string   `json:"app_name"`
	BaseURLs          []string `json:"base_urls"`
	CanRunOffline     bool     `json:"can_run_offline"`
	EGLGUID           string   `json:"egl_guid"`
	Executable        string   `json:"executable"`
	InstallPath       string   `json:"install_path"`
	InstallSize       int      `json:"install_size"`
	IsDLC             bool     `json:"is_dlc"`
	LaunchParameters  string   `json:"launch_parameters"`
	ManifestPath      string   `json:"manifest_path"`
	NeedsVerification bool     `json:"needs_verification"`
	RequiresOT        bool     `json:"requires_ot"`
	SavePath          string   `json:"save_path"`
	Title             string   `json:"title"`
	Version           string   `json:"version"`
}

type Data = map[string]Game

func FindInstallationsIn(legendaryDataPath string, launcher string) ([]*common.Installation, []error) {
	legendaryInstalledPath := filepath.Join(legendaryDataPath, "installed.json")
	if _, err := os.Stat(legendaryInstalledPath); os.IsNotExist(err) {
		return nil, []error{errors.Errorf("%s not installed", launcher)}
	}
	var legendaryData Data
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

	installs := make([]*common.Installation, 0)
	var findErrors []error

	for _, legendaryGame := range legendaryData {
		installType, version, err := common.GetGameInfo(legendaryGame.InstallPath)
		if err != nil {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  legendaryGame.InstallPath,
				Inner: err,
			})
			continue
		}

		branch, err := epic.GetEpicBranch(legendaryGame.AppName)
		if err != nil {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  legendaryGame.InstallPath,
				Inner: err,
			})
			continue
		}

		var launchPath []string
		if canLaunchLegendary {
			launchPath = []string{"legendary", "launch", legendaryGame.AppName}
		}
		installs = append(installs, &common.Installation{
			Path:       filepath.Clean(legendaryGame.InstallPath),
			Version:    version,
			Type:       installType,
			Location:   common.LocationTypeLocal,
			Branch:     branch,
			Launcher:   launcher,
			LaunchPath: launchPath,
		})
	}
	return installs, findErrors
}

func getGlobalLegendaryDataPath(xdgConfigHomeEnv string) (string, error) {
	// Should be kept in sync with
	// https://github.com/derrod/legendary/blob/master/legendary/lfs/lgndry.py#L29-L34

	if legendaryConfigPathEnv, found := os.LookupEnv("LEGENDARY_CONFIG_PATH"); found {
		return legendaryConfigPathEnv, nil
	}
	if xdgConfigHomeEnv != "" {
		return filepath.Join(xdgConfigHomeEnv, "legendary"), nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home dir")
	}
	return filepath.Join(homeDir, ".config", "legendary"), nil
}
