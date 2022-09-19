package bindings

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/SatisfactoryModManager/install_finders"
	"github.com/satisfactorymodding/SatisfactoryModManager/project_file"
	"github.com/satisfactorymodding/SatisfactoryModManager/utils"
	ficsitCli "github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type DebugInfo struct {
	ctx context.Context
}

func MakeDebugInfo() *DebugInfo {
	return &DebugInfo{}
}

func (d *DebugInfo) startup(ctx context.Context) {
	d.ctx = ctx
}

type MetadataInstallation struct {
	*install_finders.Installation
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

type Metadata struct {
	Installations        []*MetadataInstallation `json:"installations"`
	SelectedInstallation *MetadataInstallation   `json:"selectedInstallation"`
	Profiles             []*ficsitCli.Profile    `json:"profiles"`
	SelectedProfileName  *string                 `json:"selectedProfile"`
	InstalledMods        map[string]string       `json:"installedMods"`
	SMLVersion           *string                 `json:"smlVersion"`
	SMMVersion           string                  `json:"smmVersion"`
	ModsEnabled          bool                    `json:"modsEnabled"`
}

func addFactoryGameLog(writer *zip.Writer) error {
	if runtime.GOOS == "windows" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return errors.Wrap(err, "Failed to get user cache dir")
		}
		err = utils.AddFileToZip(writer, path.Join(cacheDir, "FactoryGame", "Saved", "Logs", "FactoryGame.log"), "FactoryGame.log")
		if err != nil {
			if !os.IsNotExist(err) {
				return errors.Wrap(err, "Failed to add file to zip")
			}
		}
	}
	return nil
}

func addMetadata(writer *zip.Writer) error {
	ficsitCliInstalls := BindingsInstance.FicsitCLI.GetInstallationsInfo()
	selectedFicsitCliInstall := BindingsInstance.FicsitCLI.GetSelectedInstall()
	var metadataInstalls []*MetadataInstallation
	var selectedMetadataInstall *MetadataInstallation
	for _, install := range ficsitCliInstalls {
		i := &MetadataInstallation{
			Installation: install.Info,
			Name:         fmt.Sprintf("Satisfactory %s (%s)", install.Info.Branch, install.Info.Launcher),
			Profile:      install.Installation.Profile,
		}

		metadataInstalls = append(metadataInstalls, i)

		if selectedFicsitCliInstall != nil && selectedFicsitCliInstall.Installation == install.Installation {
			selectedMetadataInstall = i
		}
	}

	ficsitCliProfileNames := BindingsInstance.FicsitCLI.GetProfiles()
	selectedMetadataProfileName := BindingsInstance.FicsitCLI.GetSelectedProfile()
	var metadataProfiles []*ficsitCli.Profile
	for _, profileName := range ficsitCliProfileNames {
		p := BindingsInstance.FicsitCLI.GetProfile(profileName)

		metadataProfiles = append(metadataProfiles, p)
	}

	lockfile, err := BindingsInstance.FicsitCLI.GetCurrentLockfile(selectedFicsitCliInstall)
	if err != nil {
		return errors.Wrap(err, "Failed to get lockfile")
	}

	metadataInstalledMods := make(map[string]string)
	var smlVersion *string
	for name, data := range *lockfile {
		if name == "SML" {
			smlVersion = &data.Version
		} else {
			metadataInstalledMods[name] = data.Version
		}
	}

	metadata := Metadata{
		SMMVersion:           project_file.Version(),
		Installations:        metadataInstalls,
		SelectedInstallation: selectedMetadataInstall,
		Profiles:             metadataProfiles,
		SelectedProfileName:  selectedMetadataProfileName,
		InstalledMods:        metadataInstalledMods,
		SMLVersion:           smlVersion,
	}

	metadataBytes, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return errors.Wrap(err, "Failed to marshal metadata")
	}

	metadataFile, err := writer.Create("metadata.json")
	if err != nil {
		return errors.Wrap(err, "Failed to create metadata file")
	}

	_, err = metadataFile.Write(metadataBytes)
	if err != nil {
		return errors.Wrap(err, "Failed to write metadata")
	}
	return nil
}

func (d *DebugInfo) generateAndSaveDebugInfo(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Failed to create file")
	}
	defer file.Close()
	writer := zip.NewWriter(file)
	defer writer.Close()

	err = addFactoryGameLog(writer)
	if err != nil {
		return errors.Wrap(err, "Failed to add FactoryGame.log to debuginfo zip")
	}

	err = addMetadata(writer)
	if err != nil {
		return errors.Wrap(err, "Failed to add metadata to debuginfo zip")
	}

	err = utils.AddFileToZip(writer, viper.GetString("log-file"), "SatisfactoryModManager.log")
	if err != nil {
		return errors.Wrap(err, "Failed to add SatisfactoryModManager.log to debuginfo zip")
	}

	return nil
}

func (d *DebugInfo) GenerateDebugInfo() bool {
	defaultFileName := fmt.Sprintf("SMMDebug-%s.zip", time.Now().UTC().Format("2006-01-02-15-04-05"))
	filename, err := wailsRuntime.SaveFileDialog(d.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: defaultFileName,
		Filters: []wailsRuntime.FileFilter{
			{
				Pattern:     "*.zip",
				DisplayName: "Zip Files (*.zip)",
			},
		},
	})
	if err != nil {
		wailsRuntime.LogErrorf(d.ctx, "Failed to get file name: %v", err)
		return false
	}
	if filename == "" {
		return false
	}

	err = d.generateAndSaveDebugInfo(filename)
	if err != nil {
		wailsRuntime.LogErrorf(d.ctx, "Failed to generate debug info: %v", err)
		return false
	}

	return true
}
