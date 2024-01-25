package bindings

import (
	"archive/zip"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	ficsitCli "github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/projectfile"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
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
			return errors.Wrap(err, "failed to get user cache dir")
		}
		err = utils.AddFileToZip(writer, path.Join(cacheDir, "FactoryGame", "Saved", "Logs", "FactoryGame.log"), "FactoryGame.log")
		if err != nil {
			if !os.IsNotExist(err) {
				return errors.Wrap(err, "failed to add file to zip")
			}
		}
	}
	return nil
}

func addMetadata(writer *zip.Writer) error {
	ficsitCliInstalls := BindingsInstance.FicsitCLI.GetInstallations()
	selectedFicsitCliInstall := BindingsInstance.FicsitCLI.GetSelectedInstall()
	metadataInstalls := make([]*MetadataInstallation, 0)
	var selectedMetadataInstall *MetadataInstallation
	for _, install := range ficsitCliInstalls {
		i := &MetadataInstallation{
			Installation: install.Info,
			Name:         fmt.Sprintf("Satisfactory %s (%s)", install.Info.Branch, install.Info.Launcher),
			Profile:      install.Installation.Profile,
		}
		i.Path = utils.RedactPath(i.Path)
		i.LaunchPath = strings.Join(i.Installation.LaunchPath, " ")

		metadataInstalls = append(metadataInstalls, i)

		if selectedFicsitCliInstall != nil && selectedFicsitCliInstall == install.Info {
			selectedMetadataInstall = i
		}
	}

	ficsitCliProfileNames := BindingsInstance.FicsitCLI.GetProfiles()
	selectedMetadataProfileName := BindingsInstance.FicsitCLI.GetSelectedProfile()
	metadataProfiles := make([]*ficsitCli.Profile, 0)
	for _, profileName := range ficsitCliProfileNames {
		p := BindingsInstance.FicsitCLI.GetProfile(profileName)

		metadataProfiles = append(metadataProfiles, p)
	}

	lockfile, err := BindingsInstance.FicsitCLI.GetSelectedInstallLockfile()
	if err != nil {
		return errors.Wrap(err, "failed to get lockfile")
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
		SMMVersion:           projectfile.Version(),
		Installations:        metadataInstalls,
		SelectedInstallation: selectedMetadataInstall,
		Profiles:             metadataProfiles,
		SelectedProfileName:  selectedMetadataProfileName,
		InstalledMods:        metadataInstalledMods,
		SMLVersion:           smlVersion,
	}

	metadataBytes, err := utils.JSONMarshal(metadata, 2)
	if err != nil {
		return errors.Wrap(err, "failed to marshal metadata")
	}

	metadataFile, err := writer.Create("metadata.json")
	if err != nil {
		return errors.Wrap(err, "failed to create metadata file")
	}

	_, err = metadataFile.Write(metadataBytes)
	if err != nil {
		return errors.Wrap(err, "failed to write metadata")
	}
	return nil
}

func (d *DebugInfo) generateAndSaveDebugInfo(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer file.Close()
	writer := zip.NewWriter(file)
	defer writer.Close()

	err = addFactoryGameLog(writer)
	if err != nil {
		return errors.Wrap(err, "failed to add FactoryGame.log to debuginfo zip")
	}

	err = addMetadata(writer)
	if err != nil {
		return errors.Wrap(err, "failed to add metadata to debuginfo zip")
	}

	err = utils.AddFileToZip(writer, viper.GetString("log-file"), "SatisfactoryModManager.log")
	if err != nil {
		return errors.Wrap(err, "failed to add SatisfactoryModManager.log to debuginfo zip")
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
		slog.Error("failed to open save dialog", slog.Any("error", err))
		return false
	}
	if filename == "" {
		return false
	}

	err = d.generateAndSaveDebugInfo(filename)
	if err != nil {
		slog.Error("failed to generate debug info", slog.Any("error", err))
		return false
	}

	return true
}
