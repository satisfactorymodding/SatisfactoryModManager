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
	"regexp"
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

func addDefaultFactoryGameLog(writer *zip.Writer) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("failed to get user cache dir: %w", err)
	}
	defaultLogPath := filepath.Join(cacheDir, "FactoryGame", "Saved", "Logs", "FactoryGame.log")
	// Default log will always be on the local disk, if it exists
	bytes, err := os.ReadFile(defaultLogPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if runtime.GOOS != "windows" {
				// On non-Windows systems the game is not running natively,
				// so the log file does not exist
				return nil
			}
			return fmt.Errorf("log does not exist")
		}
		return fmt.Errorf("failed to read default FactoryGame.log: %w", err)
	}
	err = addLogFromBytes(writer, bytes, "FactoryGame.log")
	if err != nil {
		return fmt.Errorf("failed to add default FactoryGame.log to zip: %w", err)
	}
	return nil
}

func redactLogBytes(bytes []byte) []byte {
	// Prevent leaking of Steam/Epic friend nicknames in submitted logs
	re := regexp.MustCompile(`(Added friend with nickname ').*(' on online context)`)
	return re.ReplaceAll(bytes, []byte("${1}REDACTED${2}"))
}

func addLogFromBytes(writer *zip.Writer, bytes []byte, zipFileName string) error {
	redactedBytes := redactLogBytes(bytes)

	logFile, err := writer.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("failed to create log file in zip: %w", err)
	}

	_, err = logFile.Write(redactedBytes)
	if err != nil {
		return fmt.Errorf("failed to write log file to zip: %w", err)
	}
	return nil
}

func addInstallFactoryGameLog(writer *zip.Writer, install *common.Installation) error {
	logPath := filepath.Join(install.SavedPath, "Logs", "FactoryGame.log")
	// Install-specific logs could be on remote disks
	d, err := ficsitcli.FicsitCLI.GetInstallation(install.Path).GetDisk()
	if err != nil {
		return fmt.Errorf("failed to get disk for installation: %w", err)
	}

	logExists, err := d.Exists(logPath)
	if err != nil {
		return fmt.Errorf("failed to check if log exists: %w", err)
	}
	if !logExists {
		return fmt.Errorf("log does not exist")
	}

	bytes, err := d.Read(logPath)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}
	return addLogFromBytes(writer, bytes, getLogNameForInstall(install))
}

func addFactoryGameLogs(writer *zip.Writer) {
	err := addDefaultFactoryGameLog(writer)
	if err != nil {
		slog.Warn("failed to add default FactoryGame.log to debuginfo zip", slog.Any("error", err))
	}
	for _, installMeta := range ficsitcli.FicsitCLI.GetInstallationsMetadata() {
		if installMeta.Info == nil {
			continue
		}

		err := addInstallFactoryGameLog(writer, installMeta.Info)
		if err != nil {
			slog.Warn("failed to add FactoryGame.log to debuginfo zip", slog.String("path", installMeta.Info.Path), slog.Any("error", err))
		}
	}
}

func getLogNameForInstall(install *common.Installation) string {
	hash := sha256.Sum256([]byte(install.Path))
	first8 := hex.EncodeToString(hash[:])[:8]
	return fmt.Sprintf("FactoryGame_%s_%s_%s_%s.log", first8, install.Location, install.Branch, install.Type)
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
		slog.Warn("failed to get lockfile for debuginfo", slog.Any("error", err))
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

	addFactoryGameLogs(writer)

	err = addMetadata(writer)
	if err != nil {
		slog.Warn("failed to add metadata to debuginfo zip", slog.Any("error", err))
	}

	// Add SMM log last, as it may list errors from previous steps
	err = utils.AddFileToZip(writer, viper.GetString("log-file"), "SatisfactoryModManager.log")
	if err != nil {
		return fmt.Errorf("failed to add SatisfactoryModManager.log to debuginfo zip: %w", err)
	}

	return nil
}

func (a *app) GenerateDebugInfo() (bool, error) {
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
		return false, fmt.Errorf("failed to open save dialog: %w", err)
	}
	if filename == "" {
		// user canceled the save dialog
		return false, nil
	}

	err = a.generateAndSaveDebugInfo(filename)
	if err != nil {
		slog.Error("failed to generate debug info", slog.Any("error", err))
		return false, fmt.Errorf("failed to generate debug info: %w", err)
	}

	return true, nil
}
