package ficsitcli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func (f *ficsitCLI) GetRemoteInstallations() []string {
	paths := make([]string, 0, len(f.installationMetadata))
	for _, install := range f.GetInstallations() {
		meta, ok := f.installationMetadata[install]
		if ok {
			if meta.Location != common.LocationTypeRemote {
				continue
			}
		}
		// Missing metadata means an unavailable remote installation
		paths = append(paths, install)
	}
	return paths
}

func (f *ficsitCLI) AddRemoteServer(path string) error {
	installation := f.ficsitCli.Installations.GetInstallation(path)
	if installation == nil {
		fallbackProfile := "Default"
		if f.ficsitCli.Profiles.GetProfile(fallbackProfile) == nil {
			// Pick first profile found
			for name := range f.ficsitCli.Profiles.Profiles {
				fallbackProfile = name
				break
			}
		}

		var err error
		installation, err = f.ficsitCli.Installations.AddInstallation(f.ficsitCli, path, fallbackProfile)
		if err != nil {
			return fmt.Errorf("failed to add installation: %w", err)
		}
		_ = f.ficsitCli.Installations.Save()
	}
	gameVersion, err := installation.GetGameVersion(f.ficsitCli)
	if err != nil {
		return fmt.Errorf("failed to get game version: %w", err)
	}

	platform, err := installation.GetPlatform(f.ficsitCli)
	if err != nil {
		return fmt.Errorf("failed to get platform: %w", err)
	}
	var installType common.InstallType
	switch platform.TargetName {
	case "Windows":
		installType = common.InstallTypeWindowsClient
	case "WindowsServer":
		installType = common.InstallTypeWindowsServer
	case "LinuxServer":
		installType = common.InstallTypeLinuxServer
	}

	if installType == common.InstallTypeWindowsClient {
		return fmt.Errorf("remote server is not a server installation")
	}

	branch := common.BranchEarlyAccess // TODO: Do we have a way to detect this for remote installs?

	f.installationMetadata[path] = &common.Installation{
		Path:     installation.Path,
		Type:     installType,
		Location: common.LocationTypeRemote,
		Branch:   branch,
		Version:  gameVersion,
		Launcher: f.getNextRemoteLauncherName(),
	}

	f.EmitGlobals()

	return nil
}

func (f *ficsitCLI) getNextRemoteLauncherName() string {
	existingNumbers := make(map[int]bool)
	for _, install := range f.GetRemoteInstallations() {
		metadata := f.installationMetadata[install]
		if metadata != nil {
			if strings.HasPrefix(metadata.Launcher, "Remote ") {
				num, err := strconv.Atoi(strings.TrimPrefix(metadata.Launcher, "Remote "))
				if err == nil {
					existingNumbers[num] = true
				}
			}
		}
	}
	for i := 1; ; i++ {
		if !existingNumbers[i] {
			return "Remote " + strconv.Itoa(i)
		}
	}
}

func (f *ficsitCLI) RemoveRemoteServer(path string) error {
	metadata := f.installationMetadata[path]
	if metadata == nil {
		return fmt.Errorf("installation not found")
	}
	if metadata.Location != common.LocationTypeRemote {
		return fmt.Errorf("installation is not remote")
	}
	err := f.ficsitCli.Installations.DeleteInstallation(path)
	if err != nil {
		return fmt.Errorf("failed to delete installation: %w", err)
	}
	delete(f.installationMetadata, path)
	f.EmitGlobals()
	return nil
}
