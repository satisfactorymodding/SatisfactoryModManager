package windows

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/andygrunwald/vdf"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"

	"github.com/satisfactorymodding/SatisfactoryModManager/installFinders/types"
)

func FindInstallationsSteam() ([]*types.Installation, []string, []error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return nil, nil, []error{errors.Wrap(err, "Failed to open Steam registry key")}
	}
	defer key.Close()

	steamExePath, _, err := key.GetStringValue("SteamExe")
	if err != nil {
		steamExePath = `C:\Program Files (x86)\Steam\steam.exe`
	}

	steamPath := path.Dir(steamExePath)
	steamAppsPath := path.Join(steamPath, "steamapps")
	libraryFoldersManifestPath := path.Join(steamAppsPath, "libraryfolders.vdf")

	libraryFoldersF, err := os.Open(libraryFoldersManifestPath)
	if err != nil {
		return nil, nil, []error{errors.Wrap(err, "Failed to open library folders manifest")}
	}

	parser := vdf.NewParser(libraryFoldersF)
	libraryFoldersManifest, err := parser.Parse()
	if err != nil {
		return nil, nil, []error{errors.Wrap(err, "Failed to parse library folders manifest")}
	}

	var libraryFoldersList map[string]interface{}

	if _, ok := libraryFoldersManifest["LibraryFolders"]; ok {
		libraryFoldersList = libraryFoldersManifest["LibraryFolders"].(map[string]interface{})
	} else if _, ok := libraryFoldersManifest["libraryfolders"]; ok {
		libraryFoldersList = libraryFoldersManifest["libraryfolders"].(map[string]interface{})
	} else {
		return nil, nil, []error{errors.New("Failed to find library folders in manifest")}
	}

	libraryFolders := []string{
		path.Clean(steamPath),
	}

	for key, val := range libraryFoldersList {
		if _, err := strconv.Atoi(key); err != nil {
			continue
		}

		libraryFolderData := val.(map[string]interface{})
		libraryFolders = append(libraryFolders, libraryFolderData["path"].(string))
	}

	installs := []*types.Installation{}
	invalidInstalls := []string{}
	findErrors := []error{}

	for _, libraryFolder := range libraryFolders {
		manifestPath := path.Join(libraryFolder, "steamapps", "appmanifest_526870.acf")

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
			findErrors = append(findErrors, errors.New(fmt.Sprintf("Failed to find AppState in manifest %s", manifestPath)))
			continue
		}

		fullInstallationPath := path.Join(libraryFolder, "steamapps", "common", manifest["AppState"].(map[string]interface{})["installdir"].(string))

		gameExe := path.Join(fullInstallationPath, "FactoryGame.exe")
		if _, err := os.Stat(gameExe); os.IsNotExist(err) {
			invalidInstalls = append(invalidInstalls, fullInstallationPath)
			continue
		}

		versionFilePath := path.Join(fullInstallationPath, "Engine", "Binaries", "Win64", "FactoryGame-Win64-Shipping.version")
		if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
			invalidInstalls = append(invalidInstalls, fullInstallationPath)
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
		userConfig := manifest["AppState"].(map[string]interface{})["UserConfig"].(map[string]interface{})
		betakey, ok := userConfig["betakey"]
		if !ok {
			branch = "Early Access"
		} else {
			if betakey == "experimental" {
				branch = "Experimental"
			} else {
				findErrors = append(findErrors, errors.New(fmt.Sprintf("Unknown beta key %s", betakey)))
			}
		}

		installs = append(installs, &types.Installation{
			Path:       fullInstallationPath,
			Version:    versionData.Changelist,
			Branch:     branch,
			Launcher:   "Steam",
			LaunchPath: `start "" "steam://rungameid/526870"`,
		})
	}

	return installs, invalidInstalls, findErrors
}
