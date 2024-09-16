package steam

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/andygrunwald/vdf"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

var manifests = []string{"appmanifest_526870.acf", "appmanifest_1690800.acf"}

func FindInstallationsSteam(steamPath string, launcher string, platform common.LauncherPlatform) ([]*common.Installation, []error) {
	steamAppsPath := filepath.Join(steamPath, "steamapps")
	libraryFoldersManifestPath := platform.ProcessPath(filepath.Join(steamAppsPath, "libraryfolders.vdf"))

	libraryFoldersF, err := os.Open(libraryFoldersManifestPath)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to open library folders manifest: %w", err)}
	}

	parser := vdf.NewParser(libraryFoldersF)
	libraryFoldersManifest, err := parser.Parse()
	if err != nil {
		return nil, []error{fmt.Errorf("failed to parse library folders manifest: %w", err)}
	}

	var libraryFoldersList map[string]interface{}

	if _, ok := libraryFoldersManifest["LibraryFolders"]; ok {
		libraryFoldersList = libraryFoldersManifest["LibraryFolders"].(map[string]interface{})
	} else if _, ok := libraryFoldersManifest["libraryfolders"]; ok {
		libraryFoldersList = libraryFoldersManifest["libraryfolders"].(map[string]interface{})
	} else {
		return nil, []error{fmt.Errorf("failed to find library folders in manifest")}
	}

	libraryFolders := []string{
		filepath.Clean(steamPath),
	}

	for key, val := range libraryFoldersList {
		if _, err := strconv.Atoi(key); err != nil {
			continue
		}

		libraryFolderData := val.(map[string]interface{})
		libraryFolder := libraryFolderData["path"].(string)

		found := false
		for _, existingLibraryFolder := range libraryFolders {
			if common.OsPathEqual(existingLibraryFolder, libraryFolder) {
				found = true
				break
			}
		}
		if !found {
			libraryFolders = append(libraryFolders, libraryFolder)
		}
	}

	installs := make([]*common.Installation, 0)
	var findErrors []error

	for _, libraryFolder := range libraryFolders {
		for _, manifest := range manifests {
			manifestPath := platform.ProcessPath(filepath.Join(libraryFolder, "steamapps", manifest))

			if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
				continue
			}

			manifestF, err := os.Open(manifestPath)
			if err != nil {
				findErrors = append(findErrors, fmt.Errorf("failed to open manifest file %s: %w", manifestPath, err))
				continue
			}

			parser := vdf.NewParser(manifestF)
			manifest, err := parser.Parse()
			if err != nil {
				findErrors = append(findErrors, fmt.Errorf("failed to parse manifest file %s: %w", manifestPath, err))
				continue
			}

			if _, ok := manifest["AppState"]; !ok {
				findErrors = append(findErrors, fmt.Errorf("failed to find AppState in manifest %s", manifestPath))
				continue
			}

			appState := manifest["AppState"].(map[string]interface{})

			fullInstallationPath := platform.ProcessPath(filepath.Join(libraryFolder, "steamapps", "common", appState["installdir"].(string)))

			gamePlatform := platform.Platform
			if platform.Os() != "windows" {
				// The game might be running under Proton
				// There's no appmanifest field that would specify it, but if the proton prefix exists,
				// the game is most likely running under Proton.
				gameProtonPrefix := platform.ProcessPath(filepath.Join(steamPath, "steamapps", "compatdata", appState["appid"].(string), "pfx"))
				_, err := os.Stat(gameProtonPrefix)
				if err != nil && !os.IsNotExist(err) {
					findErrors = append(findErrors, fmt.Errorf("failed to find proton prefix for game %s: %w", appState["appid"].(string), err))
					continue
				}
				if err == nil {
					gamePlatform = common.WineLauncherPlatform(gameProtonPrefix)
				}
			}

			installType, version, savedPath, err := common.GetGameInfo(fullInstallationPath, gamePlatform)
			if err != nil {
				findErrors = append(findErrors, common.InstallFindError{
					Path:  fullInstallationPath,
					Inner: err,
				})
				continue
			}

			var branch common.GameBranch
			userConfig := manifest["AppState"].(map[string]interface{})["UserConfig"].(map[string]interface{})
			betakey, ok := userConfig["betakey"]
			if !ok {
				branch = common.BranchStable
			} else {
				if betakey == "experimental" {
					branch = common.BranchExperimental
				} else {
					findErrors = append(findErrors, fmt.Errorf("unknown beta key %s", betakey))
				}
			}

			installs = append(installs, &common.Installation{
				Path:       filepath.Clean(fullInstallationPath),
				Version:    version,
				Type:       installType,
				Location:   common.LocationTypeLocal,
				Branch:     branch,
				Launcher:   launcher,
				LaunchPath: platform.LauncherCommand(`steam://rungameid/526870`),
				// pass wine platform if necessary, as platform here is going to be native
				SavedPath: savedPath,
			})
		}
	}

	return installs, findErrors
}
