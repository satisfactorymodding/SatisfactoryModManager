package legendary

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"

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

func FindInstallationsIn(legendaryDataPath string, launcher string, knownPrefixes map[string]string, platform common.LauncherPlatform) ([]*common.Installation, []error) {
	legendaryInstalledPath := filepath.Join(legendaryDataPath, "installed.json")
	if _, err := os.Stat(legendaryInstalledPath); os.IsNotExist(err) {
		return nil, []error{fmt.Errorf("%s not installed", launcher)}
	}
	var legendaryData Data
	legendaryDataFile, err := os.ReadFile(legendaryInstalledPath)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to read legendary installed.json: %w", err)}
	}
	err = json.Unmarshal(legendaryDataFile, &legendaryData)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to parse legendary installed.json output: %w", err)}
	}

	installs := make([]*common.Installation, 0)
	var findErrors []error

	for _, legendaryGame := range legendaryData {
		installLocation := filepath.Clean(legendaryGame.InstallPath)

		gamePlatform := platform.Platform

		if platform.Os() != "windows" {
			if knownPrefix, found := knownPrefixes[legendaryGame.AppName]; found {
				gamePlatform = common.WineLauncherPlatform(knownPrefix)
			} else {
				prefix, err := getLegendaryWinePrefix(legendaryDataPath, legendaryGame.AppName, platform)
				if err != nil {
					findErrors = append(findErrors, fmt.Errorf("failed to get wine prefix for %s: %w", legendaryGame.AppName, err))
					continue
				}
				if prefix != "" {
					gamePlatform = common.WineLauncherPlatform(prefix)
				}
			}
		}

		installType, version, savedPath, err := common.GetGameInfo(installLocation, gamePlatform)
		if err != nil {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
				Inner: err,
			})
			continue
		}

		branch, err := epic.GetEpicBranch(legendaryGame.AppName)
		if err != nil {
			findErrors = append(findErrors, common.InstallFindError{
				Path:  installLocation,
				Inner: err,
			})
			continue
		}

		installs = append(installs, &common.Installation{
			Path:       installLocation,
			Version:    version,
			Type:       installType,
			Location:   common.LocationTypeLocal,
			Branch:     branch,
			Launcher:   launcher,
			LaunchPath: platform.LauncherCommand(legendaryGame.AppName),
			SavedPath:  savedPath,
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
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}
	return filepath.Join(homeDir, ".config", "legendary"), nil
}

func getLegendaryWinePrefix(legendaryDataPath string, appName string, platform common.Platform) (string, error) {
	// Should be kept in sync with
	// https://github.com/derrod/legendary/blob/master/legendary/core.py#L591

	config, err := ini.Load(filepath.Join(legendaryDataPath, "config.ini"))
	if err != nil {
		return "", fmt.Errorf("failed to load legendary config.ini: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}

	prefix := ""

	prefix = stringOrFallback(config.Section("default.env").Key("WINEPREFIX").String(), prefix)
	prefix = stringOrFallback(config.Section(fmt.Sprintf("%s.env", appName)).Key("WINEPREFIX").String(), prefix)

	if platform.Os() == "darwin" {
		cxBottle := "Legendary"
		cxBottle = stringOrFallback(config.Section("default").Key("crossover_bottle").String(), cxBottle)
		cxBottle = stringOrFallback(config.Section(appName).Key("crossover_bottle").String(), cxBottle)

		bottlePath := filepath.Join(homeDir, "Library", "Application Support", "CrossOver", "Bottles", cxBottle)
		if _, err := os.Stat(bottlePath); err == nil {
			prefix = stringOrFallback(bottlePath, prefix)
		}
	}

	prefix = stringOrFallback(config.Section(appName).Key("wine_prefix").String(), prefix)

	prefix = stringOrFallback(filepath.Join(homeDir, ".wine"), prefix)

	return prefix, nil
}

func stringOrFallback(a, fallback string) string {
	if a == "" {
		return fallback
	}
	return a
}
