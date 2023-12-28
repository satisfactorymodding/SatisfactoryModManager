package ficsitcli

import (
	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func (f *FicsitCLI) GetRemoteInstallations() []*common.Installation {
	remoteInstallations := []*common.Installation{}
	for _, install := range f.installations {
		if install.Info.Location == common.LocationTypeRemote {
			remoteInstallations = append(remoteInstallations, install.Info)
		}
	}
	return remoteInstallations
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

	branch := common.BranchEarlyAccess // TODO: Do we have a way to detect this for remote installs?

	f.installations = append(f.installations, &InstallationInfo{
		Installation: installation,
		Info: &common.Installation{
			Path:     installation.Path,
			Type:     installType,
			Location: common.LocationTypeRemote,
			Branch:   branch,
			Version:  gameVersion,
			Launcher: "Remote",
		},
	})

	f.EmitGlobals()

	return nil
}

func (f *FicsitCLI) RemoveRemoteServer(path string) error {
	for _, install := range f.installations {
		if install.Info.Path == path && install.Info.Location != common.LocationTypeRemote {
			return errors.New("installation is not remote")
		}
	}
	err := f.ficsitCli.Installations.DeleteInstallation(path)
	if err != nil {
		return errors.Wrap(err, "failed to delete installation")
	}
	for i, install := range f.installations {
		if install.Info.Path == path {
			f.installations = append(f.installations[:i], f.installations[i+1:]...)
			break
		}
	}
	f.EmitGlobals()
	return nil
}
