package app

import (
	"archive/zip"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	ficsitCli "github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	appCommon "github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/ficsitcli"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

type MetadataInstallation struct {
	*common.Installation
	LaunchPath string `json:"launchPath"`
	Name       string `json:"name"`
	Profile    string `json:"profile"`
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
			return fmt.Errorf("failed to get user cache dir: %w", err)
		}
		err = utils.AddFileToZip(writer, filepath.Join(cacheDir, "FactoryGame", "Saved", "Logs", "FactoryGame.log"), "FactoryGame.log")
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to add file to zip: %w", err)
			}
		}
	}
	return nil
}

func addMetadata(writer *zip.Writer) error {
	installs := ficsitcli.FicsitCLI.GetInstallations()
	selectedInstallInstance := ficsitcli.FicsitCLI.GetSelectedInstall()
	metadataInstalls := make([]*MetadataInstallation, 0)
	var selectedMetadataInstall *MetadataInstallation
	for _, install := range installs {
		metadata := ficsitcli.FicsitCLI.GetInstallationsMetadata()[install]
		if metadata.Info == nil {
			slog.Warn("failed to get metadata for installation", utils.SlogPath("path", install))
			continue
		}
		i := &MetadataInstallation{
			Installation: metadata.Info,
			Name:         fmt.Sprintf("Satisfactory %s (%s)", metadata.Info.Branch, metadata.Info.Branch),
			Profile:      ficsitcli.FicsitCLI.GetInstallation(install).Profile,
		}
		i.Path = utils.RedactPath(i.Path)
		i.LaunchPath = strings.Join(i.Installation.LaunchPath, " ")

		metadataInstalls = append(metadataInstalls, i)

		if selectedInstallInstance != nil && selectedInstallInstance.Path == install {
			selectedMetadataInstall = i
		}
	}

	ficsitCliProfileNames := ficsitcli.FicsitCLI.GetProfiles()
	selectedMetadataProfileName := ficsitcli.FicsitCLI.GetSelectedProfile()
	metadataProfiles := make([]*ficsitCli.Profile, 0)
	for _, profileName := range ficsitCliProfileNames {
		p := ficsitcli.FicsitCLI.GetProfile(profileName)

		metadataProfiles = append(metadataProfiles, p)
	}

	lockfile, err := ficsitcli.FicsitCLI.GetSelectedInstallLockfile()
	if err != nil {
		return fmt.Errorf("failed to get lockfile: %w", err)
	}

	metadataInstalledMods := make(map[string]string)
	var smlVersion *string

	if lockfile != nil {
		for name, data := range lockfile.Mods {
			if name == "SML" {
				smlVersion = &data.Version
			} else {
				metadataInstalledMods[name] = data.Version
			}
		}
	}

	metadata := Metadata{
		SMMVersion:           viper.GetString("version"),
		Installations:        metadataInstalls,
		SelectedInstallation: selectedMetadataInstall,
		Profiles:             metadataProfiles,
		SelectedProfileName:  selectedMetadataProfileName,
		InstalledMods:        metadataInstalledMods,
		SMLVersion:           smlVersion,
	}

	metadataBytes, err := utils.JSONMarshal(metadata, 2)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataFile, err := writer.Create("metadata.json")
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}

	_, err = metadataFile.Write(metadataBytes)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}
	return nil
}

func (a *app) generateAndSaveDebugInfo(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	writer := zip.NewWriter(file)
	defer writer.Close()

	err = addFactoryGameLog(writer)
	if err != nil {
		return fmt.Errorf("failed to add FactoryGame.log to debuginfo zip: %w", err)
	}

	err = addMetadata(writer)
	if err != nil {
		return fmt.Errorf("failed to add metadata to debuginfo zip: %w", err)
	}

	err = utils.AddFileToZip(writer, viper.GetString("log-file"), "SatisfactoryModManager.log")
	if err != nil {
		return fmt.Errorf("failed to add SatisfactoryModManager.log to debuginfo zip: %w", err)
	}

	return nil
}

func (a *app) GenerateDebugInfo() bool {
	defaultFileName := fmt.Sprintf("SMMDebug-%s.zip", time.Now().UTC().Format("2006-01-02-15-04-05"))
	filename, err := wailsRuntime.SaveFileDialog(appCommon.AppContext, wailsRuntime.SaveDialogOptions{
		DefaultFilename: defaultFileName,
		Filters: []wailsRuntime.FileFilter{
			{
				Pattern:     "*.zip",
				DisplayName: "Zip Files (*.zip)",
			},
		},
	})
	if err != nil {
		slog.Error("failed to open save dialog", slog.Any("error", err))
		return false
	}
	if filename == "" {
		return false
	}

	err = a.generateAndSaveDebugInfo(filename)
	if err != nil {
		slog.Error("failed to generate debug info", slog.Any("error", err))
		return false
	}

	return true
}
