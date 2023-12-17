package installfinders

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/andygrunwald/vdf"
	"github.com/pkg/errors"
)

func findInstallationsSteam(steamPath string, launcher string, executable []string) ([]*Installation, []error) {
	steamAppsPath := filepath.Join(steamPath, "steamapps")
	libraryFoldersManifestPath := filepath.Join(steamAppsPath, "libraryfolders.vdf")

	libraryFoldersF, err := os.Open(libraryFoldersManifestPath)
	if err != nil {
		return nil, []error{errors.Wrap(err, "Failed to open library folders manifest")}
	}

	parser := vdf.NewParser(libraryFoldersF)
	libraryFoldersManifest, err := parser.Parse()
	if err != nil {
		return nil, []error{errors.Wrap(err, "Failed to parse library folders manifest")}
	}

	var libraryFoldersList map[string]interface{}

	if _, ok := libraryFoldersManifest["LibraryFolders"]; ok {
		libraryFoldersList = libraryFoldersManifest["LibraryFolders"].(map[string]interface{})
	} else if _, ok := libraryFoldersManifest["libraryfolders"]; ok {
		libraryFoldersList = libraryFoldersManifest["libraryfolders"].(map[string]interface{})
	} else {
		return nil, []error{errors.New("Failed to find library folders in manifest")}
	}

	libraryFolders := []string{
		path.Clean(steamPath),
	}

	for key, val := range libraryFoldersList {
		if _, err := strconv.Atoi(key); err != nil {
			continue
		}

		libraryFolderData := val.(map[string]interface{})
		libraryFolder := libraryFolderData["path"].(string)

		found := false
		for _, existingLibraryFolder := range libraryFolders {
			if OsPathEqual(existingLibraryFolder, libraryFolder) {
				found = true
				break
			}
		}
		if !found {
			libraryFolders = append(libraryFolders, libraryFolder)
		}
	}

	installs := make([]*Installation, 0)
	var findErrors []error

	for _, libraryFolder := range libraryFolders {
		manifestPath := filepath.Join(libraryFolder, "steamapps", "appmanifest_526870.acf")

		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		manifestF, err := os.Open(manifestPath)
		if err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to open manifest file %s", manifestPath))
			continue
		}

		parser := vdf.NewParser(manifestF)
		manifest, err := parser.Parse()
		if err != nil {
			findErrors = append(findErrors, errors.Wrapf(err, "Failed to parse manifest file %s", manifestPath))
			continue
		}

		if _, ok := manifest["AppState"]; !ok {
			findErrors = append(findErrors, errors.Errorf("Failed to find AppState in manifest %s", manifestPath))
			continue
		}

		fullInstallationPath := filepath.Join(libraryFolder, "steamapps", "common", manifest["AppState"].(map[string]interface{})["installdir"].(string))

		gameExe := filepath.Join(fullInstallationPath, "FactoryGame.exe")
		if _, err := os.Stat(gameExe); os.IsNotExist(err) {
			findErrors = append(findErrors, InstallFindError{
				Path:  fullInstallationPath,
				Inner: errors.Wrap(err, "Missing game executable"),
			})
			continue
		}

		versionFilePath := filepath.Join(fullInstallationPath, "Engine", "Binaries", "Win64", "FactoryGame-Win64-Shipping.version")
		if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
			findErrors = append(findErrors, InstallFindError{
				Path:  fullInstallationPath,
				Inner: errors.Wrap(err, "Missing game version file"),
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
		userConfig := manifest["AppState"].(map[string]interface{})["UserConfig"].(map[string]interface{})
		betakey, ok := userConfig["betakey"]
		if !ok {
			branch = BranchEarlyAccess
		} else {
			if betakey == "experimental" {
				branch = BranchExperimental
			} else {
				findErrors = append(findErrors, errors.Errorf("Unknown beta key %s", betakey))
			}
		}

		installs = append(installs, &Installation{
			Path:     fullInstallationPath,
			Version:  versionData.Changelist,
			Type:     InstallTypeWindowsClient,
			Location: LocationTypeLocal,
			Branch:   branch,
			Launcher: launcher,
			LaunchPath: append(
				executable,
				`steam://rungameid/526870`,
			),
		})
	}

	return installs, findErrors
}
