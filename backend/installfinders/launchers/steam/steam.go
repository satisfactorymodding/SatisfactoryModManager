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

func FindInstallationsSteam(steamPath string, launcher string, launchPath func(steamApp string) []string, processPath func(path string) string) ([]*common.Installation, []error) {
	if launchPath == nil {
		launchPath = func(appName string) []string { return nil }
	}

	if processPath == nil {
		processPath = func(path string) string { return path }
	}

	steamAppsPath := filepath.Join(steamPath, "steamapps")
	libraryFoldersManifestPath := processPath(filepath.Join(steamAppsPath, "libraryfolders.vdf"))

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
			manifestPath := processPath(filepath.Join(libraryFolder, "steamapps", manifest))

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

			fullInstallationPath := processPath(filepath.Join(libraryFolder, "steamapps", "common", manifest["AppState"].(map[string]interface{})["installdir"].(string)))

			installType, version, err := common.GetGameInfo(fullInstallationPath)
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
				branch = common.BranchEarlyAccess
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
				LaunchPath: launchPath(`steam://rungameid/526870`),
			})
		}
	}

	return installs, findErrors
}
