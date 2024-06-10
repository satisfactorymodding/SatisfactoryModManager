package app

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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
	Log        string `json:"log"`
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

func addFactoryGameLogs(writer *zip.Writer) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("failed to get user cache dir: %w", err)
	}
	err = utils.AddFileToZip(writer, filepath.Join(cacheDir, "FactoryGame", "Saved", "Logs", "FactoryGame.log"), "FactoryGame.log")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to add file to zip: %w", err)
		}
	}
	for _, meta := range ficsitcli.FicsitCLI.GetInstallationsMetadata() {
		if meta.Info == nil {
			continue
		}

		logPath := filepath.Join(meta.Info.SavedPath, "Logs", "FactoryGame.log")
		d, err := ficsitcli.FicsitCLI.GetInstallation(meta.Info.Path).GetDisk()
		if err != nil {
			slog.Warn("failed to get disk for installation", slog.String("path", meta.Info.Path), slog.Any("error", err))
			continue
		}
		logExists, err := d.Exists(logPath)
		if err != nil {
			slog.Warn("failed to check if log exists", slog.String("path", logPath), slog.Any("error", err))
			continue
		}
		if !logExists {
			continue
		}
		bytes, err := d.Read(logPath)
		if err != nil {
			slog.Warn("failed to read log file", slog.String("path", logPath), slog.Any("error", err))
			continue
		}
		logFile, err := writer.Create(getLogNameForInstall(meta.Info))
		if err != nil {
			slog.Warn("failed to create log file in zip", slog.Any("error", err))
			continue
		}
		_, err = logFile.Write(bytes)
		if err != nil {
			slog.Warn("failed to write log file to zip", slog.Any("error", err))
			continue
		}
	}
	return nil
}

func getLogNameForInstall(install *common.Installation) string {
	hash := sha256.Sum256([]byte(install.Path))
	return fmt.Sprintf("FactoryGame_%s.log", hex.EncodeToString(hash[:])[:8])
}

func addMetadata(writer *zip.Writer) error {
	installs := ficsitcli.FicsitCLI.GetInstallations()
	selectedInstallInstance := ficsitcli.FicsitCLI.GetSelectedInstall()
	metadataInstalls := make([]*MetadataInstallation, 0)
	var selectedMetadataInstall *MetadataInstallation
	for _, install := range installs {
		metadata := ficsitcli.FicsitCLI.GetInstallationsMetadata()[install]
		if metadata.Info == nil {
			slog.Warn("failed to get metadata for installation", slog.String("path", install))
			continue
		}
		i := &MetadataInstallation{
			Installation: metadata.Info,
			Name:         fmt.Sprintf("Satisfactory %s (%s)", metadata.Info.Branch, metadata.Info.Launcher),
			Profile:      ficsitcli.FicsitCLI.GetInstallation(install).Profile,
			Log:          getLogNameForInstall(metadata.Info),
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

	err = addFactoryGameLogs(writer)
	if err != nil {
		return fmt.Errorf("failed to add FactoryGame.log to debuginfo zip: %w", err)
	}

	err = addMetadata(writer)
	if err != nil {
		return fmt.Errorf("failed to add metadata to debuginfo zip: %w", err)
	}

	// Add SMM log last, as it may list errors from previous steps
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
