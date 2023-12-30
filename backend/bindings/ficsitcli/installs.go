package ficsitcli

import (
	"log/slog"
	"net/url"
	"os/exec"
	"sort"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	resolver "github.com/satisfactorymodding/ficsit-resolver"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

func (f *FicsitCLI) initInstallations() error {
	err := f.initLocalInstallations()
	if err != nil {
		return errors.Wrap(err, "failed to initialize found installations")
	}

	err = f.initRemoteServerInstallations()
	if err != nil {
		return errors.Wrap(err, "failed to initialize remote server installations")
	}

	sort.Slice(f.installations, func(i, j int) bool {
		if f.installations[i].Info.Launcher != f.installations[j].Info.Launcher {
			return f.installations[i].Info.Launcher < f.installations[j].Info.Launcher
		}
		return f.installations[i].Info.Branch < f.installations[j].Info.Branch
	})

	if len(f.installations) > 0 {
		f.selectedInstallation = f.installations[0]
	}

	if f.ficsitCli.Installations.SelectedInstallation != "" {
		for _, install := range f.installations {
			if install.Info.Path == f.ficsitCli.Installations.SelectedInstallation {
				f.selectedInstallation = install
				break
			}
		}
	}

	return nil
}

func (f *FicsitCLI) initLocalInstallations() error {
	installs, findErrors := installfinders.FindInstallations()

	f.installFindErrors = findErrors
	f.installations = []*InstallationInfo{}

	fallbackProfile := "Default"
	if f.ficsitCli.Profiles.GetProfile(fallbackProfile) == nil {
		// Pick first profile found
		for name := range f.ficsitCli.Profiles.Profiles {
			fallbackProfile = name
			break
		}
	}

	createdNewInstalls := false
	for _, install := range installs {
		ficsitCliInstall := f.ficsitCli.Installations.GetInstallation(install.Path)
		if ficsitCliInstall == nil {
			var err error
			ficsitCliInstall, err = f.ficsitCli.Installations.AddInstallation(f.ficsitCli, install.Path, fallbackProfile)
			if err != nil {
				return errors.Wrap(err, "failed to add installation")
			}
			createdNewInstalls = true
		}
		f.installations = append(f.installations, &InstallationInfo{
			Installation: ficsitCliInstall,
			Info:         install,
		})
	}

	if createdNewInstalls {
		err := f.ficsitCli.Installations.Save()
		if err != nil {
			return errors.Wrap(err, "failed to save installations")
		}
	}
	return nil
}

func (f *FicsitCLI) initRemoteServerInstallations() error {
	for _, installation := range f.ficsitCli.Installations.Installations {
		err := f.checkAndAddExistingRemote(installation)
		if err != nil {
			slog.Warn("failed to check and add existing remote", slog.Any("error", err), utils.SlogPath("path", installation.Path))
		}
	}
	return nil
}

func (f *FicsitCLI) checkAndAddExistingRemote(installation *cli.Installation) error {
	slog.Debug("checking whether installation is remote", utils.SlogPath("path", installation.Path))
	parsed, err := url.Parse(installation.Path)
	if err != nil {
		return errors.Wrap(err, "failed to parse installation path")
	}
	if parsed.Scheme == "ftp" || parsed.Scheme == "sftp" {
		// It is not a local installation
		if err := f.AddRemoteServer(installation.Path); err != nil {
			return errors.Wrap(err, "failed to add remote server")
		}
	}
	return nil
}

func (f *FicsitCLI) GetInstallations() []*InstallationInfo {
	return f.installations
}

func (f *FicsitCLI) GetInstallationsInfo() []*common.Installation {
	result := []*common.Installation{}
	for _, install := range f.installations {
		result = append(result, install.Info)
	}
	return result
}

func (f *FicsitCLI) GetInvalidInstalls() []string {
	result := []string{}
	for _, err := range f.installFindErrors {
		var installFindErr common.InstallFindError
		if ok := errors.As(err, &installFindErr); ok {
			result = append(result, installFindErr.Path)
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
	l := slog.With(slog.String("task", "selectInstall"), utils.SlogPath("path", path))
	if f.selectedInstallation != nil && f.selectedInstallation.Info.Path == path {
		return nil
	}
	installation := f.GetInstallation(path)
	if installation == nil {
		l.Error("failed to find installation")
		return errors.Errorf("installation %s not found", path)
	}
	f.selectedInstallation = installation

	f.ficsitCli.Installations.SelectedInstallation = path
	_ = f.ficsitCli.Installations.Save()

	f.EmitGlobals()

	f.progress = &Progress{
		Item:     "__select_install__",
		Message:  "Validating install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(f.selectedInstallation, "__select_install__")

	if installErr != nil {
		l.Error("failed to validate install", slog.Any("error", installErr))
		return errors.Wrap(installErr, "Failed to validate install")
	}
	return nil
}

func (f *FicsitCLI) GetSelectedInstall() *common.Installation {
	if f.selectedInstallation == nil {
		return nil
	}
	return f.selectedInstallation.Info
}

func (f *FicsitCLI) SetModsEnabled(enabled bool) error {
	if f.selectedInstallation == nil {
		slog.Error("no installation selected")
		return errors.New("no installation selected")
	}
	l := slog.With(slog.String("task", "setModsEnabled"), slog.Bool("enabled", enabled), utils.SlogPath("install", f.selectedInstallation.Info.Path))

	var message string
	if enabled {
		message = "Enabling mods"
	} else {
		message = "Disabling mods"
	}

	f.selectedInstallation.Installation.Vanilla = !enabled
	_ = f.ficsitCli.Installations.Save()

	f.EmitGlobals()

	f.progress = &Progress{
		Item:     "__toggle_mods__",
		Message:  message,
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(f.selectedInstallation, "__toggle_mods__")

	if installErr != nil {
		l.Error("failed to validate install", slog.Any("error", installErr))
		return errors.Wrap(installErr, "failed to validate install")
	}

	return nil
}

func (f *FicsitCLI) GetModsEnabled() bool {
	return f.selectedInstallation == nil || !f.selectedInstallation.Installation.Vanilla
}

func (f *FicsitCLI) GetSelectedInstallProfileMods() map[string]cli.ProfileMod {
	if f.selectedInstallation == nil {
		return make(map[string]cli.ProfileMod)
	}
	profile := f.GetProfile(f.selectedInstallation.Installation.Profile)
	return profile.Mods
}

func (f *FicsitCLI) GetSelectedInstallLockfileMods() (map[string]resolver.LockedMod, error) {
	if f.selectedInstallation == nil {
		return make(map[string]resolver.LockedMod), nil
	}
	lockfile, err := f.selectedInstallation.Installation.LockFile(f.ficsitCli)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get lockfile")
	}
	if lockfile == nil {
		return make(map[string]resolver.LockedMod), nil
	}
	return lockfile.Mods, nil
}

func (f *FicsitCLI) GetSelectedInstallLockfile() (*resolver.LockFile, error) {
	if f.selectedInstallation == nil {
		return nil, nil
	}
	lockfile, err := f.selectedInstallation.Installation.LockFile(f.ficsitCli)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get lockfile")
	}
	return lockfile, nil
}

func (f *FicsitCLI) LaunchGame() {
	if f.selectedInstallation == nil {
		slog.Error("no installation selected")
		return
	}
	cmd := exec.Command(f.selectedInstallation.Info.LaunchPath[0], f.selectedInstallation.Info.LaunchPath[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("failed to launch game", slog.Any("error", err), slog.String("cmd", cmd.String()), slog.String("output", string(out)))
		return
	}
}
