package common

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type installType struct {
	versionPath string
	executable  string
	installType InstallType
}

var gameInfo = []installType{
	{
		executable:  "FactoryServer.sh",
		versionPath: filepath.Join("Engine", "Binaries", "Linux", "UnrealServer-Linux-Shipping.version"),
		installType: InstallTypeLinuxServer,
	},
	{
		executable:  "FactoryServer.exe",
		versionPath: filepath.Join("Engine", "Binaries", "Win64", "UnrealServer-Win64-Shipping.version"),
		installType: InstallTypeWindowsServer,
	},
	{
		executable:  "FactoryGame.exe",
		versionPath: filepath.Join("Engine", "Binaries", "Win64", "FactoryGame-Win64-Shipping.version"),
		installType: InstallTypeWindowsClient,
	},
	// Update 9 stuff below
	{
		executable:  "FactoryServer.sh",
		versionPath: filepath.Join("Engine", "Binaries", "Linux", "FactoryServer-Linux-Shipping.version"),
		installType: InstallTypeLinuxServer,
	},
	{
		executable:  "FactoryServer.exe",
		versionPath: filepath.Join("Engine", "Binaries", "Win64", "FactoryServer-Win64-Shipping.version"),
		installType: InstallTypeWindowsServer,
	},
	// Update 1.0 stuff
	{
		executable:  "FactoryGameSteam.exe",
		versionPath: filepath.Join("Engine", "Binaries", "Win64", "FactoryGameSteam-Win64-Shipping.version"),
		installType: InstallTypeWindowsClient,
	},
	{
		executable:  "FactoryGameEGS.exe",
		versionPath: filepath.Join("Engine", "Binaries", "Win64", "FactoryGameEGS-Win64-Shipping.version"),
		installType: InstallTypeWindowsClient,
	},
}

type GameVersionFile struct {
	MajorVersion         int    `json:"MajorVersion"`
	MinorVersion         int    `json:"MinorVersion"`
	PatchVersion         int    `json:"PatchVersion"`
	Changelist           int    `json:"Changelist"`
	CompatibleChangelist int    `json:"CompatibleChangelist"`
	IsLicenseeVersion    int    `json:"IsLicenseeVersion"`
	IsPromotedBuild      int    `json:"IsPromotedBuild"`
	BranchName           string `json:"BranchName"`
	BuildID              string `json:"BuildId"`
}

func GetGameInfo(path string, platform Platform) (InstallType, int, string, error) {
	for _, info := range gameInfo {
		executablePath := filepath.Join(path, info.executable)
		if _, err := os.Stat(executablePath); os.IsNotExist(err) {
			slog.Debug("game not of type", slog.String("path", executablePath), slog.String("type", string(info.installType)))
			continue
		}

		versionFilePath := filepath.Join(path, info.versionPath)
		if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
			slog.Debug("game not of type", slog.String("path", executablePath), slog.String("type", string(info.installType)))
			continue
		}

		versionFile, err := os.ReadFile(versionFilePath)
		if err != nil {
			return InstallTypeWindowsClient, 0, "", fmt.Errorf("failed to read version file %s: %w", versionFilePath, err)
		}

		var versionData GameVersionFile
		if err := json.Unmarshal(versionFile, &versionData); err != nil {
			return InstallTypeWindowsClient, 0, "", fmt.Errorf("failed to parse version file %s: %w", versionFilePath, err)
		}

		return info.installType, versionData.Changelist, getGameSavedDir(path, info.installType, platform), nil
	}
	return InstallTypeWindowsClient, 0, "", fmt.Errorf("failed to get game info")
}

func getGameSavedDir(gamePath string, install InstallType, platform Platform) string {
	if install == InstallTypeWindowsClient {
		cacheDir, err := platform.CacheDir()
		if err != nil {
			slog.Error("failed to get cache dir", slog.Any("error", err))
			return ""
		}

		return filepath.Join(cacheDir, "FactoryGame", "Saved")
	}
	return filepath.Join(gamePath, "FactoryGame", "Saved")
}
