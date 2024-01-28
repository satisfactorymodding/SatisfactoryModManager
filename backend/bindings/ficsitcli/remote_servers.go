package ficsitcli

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func (f *FicsitCLI) GetRemoteInstallations() []string {
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

func (f *FicsitCLI) AddRemoteServer(path string) error {
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
			return errors.Wrap(err, "failed to add installation")
		}
		_ = f.ficsitCli.Installations.Save()
	}
	gameVersion, err := installation.GetGameVersion(f.ficsitCli)
	if err != nil {
		return errors.Wrap(err, "failed to get game version")
	}

	platform, err := installation.GetPlatform(f.ficsitCli)
	if err != nil {
		return errors.Wrap(err, "failed to get platform")
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
		return errors.New("remote server is not a server installation")
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

func (f *FicsitCLI) getNextRemoteLauncherName() string {
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

func (f *FicsitCLI) RemoveRemoteServer(path string) error {
	metadata := f.installationMetadata[path]
	if metadata == nil {
		return errors.New("installation not found")
	}
	if metadata.Location != common.LocationTypeRemote {
		return errors.New("installation is not remote")
	}
	err := f.ficsitCli.Installations.DeleteInstallation(path)
	if err != nil {
		return errors.Wrap(err, "failed to delete installation")
	}
	delete(f.installationMetadata, path)
	f.EmitGlobals()
	return nil
}
