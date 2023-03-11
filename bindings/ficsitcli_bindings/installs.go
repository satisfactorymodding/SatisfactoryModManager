package ficsitcli_bindings

import (
	"os/exec"
	"sort"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/SatisfactoryModManager/install_finders"
	"github.com/satisfactorymodding/SatisfactoryModManager/settings"
	"github.com/satisfactorymodding/ficsit-cli/cli"
)

func (f *FicsitCLI) initInstallations() error {
	installs, findErrors := install_finders.FindInstallations()

	f.installFindErrors = findErrors
	f.installations = []*InstallationInfo{}
	f.ficsitCli.Installations.Installations = []*cli.Installation{}

	for _, install := range installs {
		ficsitCliInstall, err := f.ficsitCli.Installations.AddInstallation(f.ficsitCli, install.Path, "Default")
		if err != nil {
			return errors.Wrap(err, "failed to add installation")
		}
		f.installations = append(f.installations, &InstallationInfo{
			Installation: ficsitCliInstall,
			Info:         install,
		})
	}

	sort.Slice(f.installations, func(i, j int) bool {
		if f.installations[i].Info.Launcher != f.installations[j].Info.Launcher {
			return f.installations[i].Info.Launcher < f.installations[j].Info.Launcher
		}
		return f.installations[i].Info.Branch < f.installations[j].Info.Branch
	})

	for _, install := range f.installations {
		savedSelectedProfile, ok := settings.Settings.SelectedProfile[install.Info.Path]
		if ok {
			err := install.Installation.SetProfile(f.ficsitCli, savedSelectedProfile)
			if err != nil {
				return errors.Wrap(err, "failed to set profile")
			}
		}
	}

	if len(f.installations) > 0 {
		f.selectedInstallation = f.installations[0]
	}

	savedSelectedInstall := settings.Settings.SelectedInstall
	if savedSelectedInstall != "" {
		for _, install := range f.installations {
			if install.Info.Path == savedSelectedInstall {
				f.selectedInstallation = install
				break
			}
		}
	}

	return nil
}

func (f *FicsitCLI) GetInstallationsInfo() []*InstallationInfo {
	return f.installations
}

func (f *FicsitCLI) GetInvalidInstalls() []string {
	result := []string{}
	for _, err := range f.installFindErrors {
		if casted, ok := err.(install_finders.InstallFindError); ok {
			result = append(result, casted.Path)
		}
	}
	return result
}

func (f *FicsitCLI) GetInstallation(path string) *InstallationInfo {
	for _, install := range f.installations {
		if install.Info.Path == path {
			return install
		}
	}

	return nil
}

func (f *FicsitCLI) SelectInstall(path string) error {
	if f.selectedInstallation != nil && f.selectedInstallation.Info.Path == path {
		return nil
	}
	installation := f.GetInstallation(path)
	if installation == nil {
		log.Error().Str("path", path).Msg("Failed to find installation")
		return errors.New("Installation \"" + path + "\" not found")
	}
	f.selectedInstallation = installation

	f.progress = &Progress{
		Item:     "__select_install__",
		Message:  "Validating install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(installation, "__select_install__")

	if installErr != nil {
		log.Error().Err(installErr).Str("install", installation.Info.Path).Msg("Failed to validate install")
		return errors.Wrap(installErr, "Failed to validate install")
	}

	settings.Settings.SelectedInstall = installation.Info.Path
	settings.SaveSettings()
	return nil
}

func (f *FicsitCLI) GetSelectedInstall() *InstallationInfo {
	return f.selectedInstallation
}

func (f *FicsitCLI) GetLockFile(installation *InstallationInfo) (*cli.LockFile, error) {
	return installation.Installation.LockFile(f.ficsitCli)
}

func (f *FicsitCLI) LaunchGame() {
	if f.selectedInstallation == nil {
		log.Error().Msg("No installation selected")
		return
	}
	cmd := exec.Command(f.selectedInstallation.Info.LaunchPath[0], f.selectedInstallation.Info.LaunchPath[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Str("cmd", cmd.String()).Str("output", string(out)).Msg("Failed to launch game")
		return
	}
}
