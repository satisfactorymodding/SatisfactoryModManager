package ficsitcli

import (
	"fmt"
	"log/slog"
	"os/exec"
	"slices"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	resolver "github.com/satisfactorymodding/ficsit-resolver"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

func (f *FicsitCLI) initInstallations() error {
	err := f.initLocalInstallationsMetadata()
	if err != nil {
		return errors.Wrap(err, "failed to initialize found installations")
	}

	err = f.initRemoteServerInstallationsMetadata()
	if err != nil {
		return errors.Wrap(err, "failed to initialize remote server installations")
	}

	filteredInstalls := f.GetInstallations()
	if len(filteredInstalls) > 0 {
		if !slices.Contains(filteredInstalls, f.ficsitCli.Installations.SelectedInstallation) {
			f.ficsitCli.Installations.SelectedInstallation = filteredInstalls[0]
			_ = f.ficsitCli.Installations.Save()
		}
	}

	return nil
}

func (f *FicsitCLI) initLocalInstallationsMetadata() error {
	installs, findErrors := installfinders.FindInstallations()

	f.installFindErrors = findErrors
	f.installationMetadata = make(map[string]*common.Installation)

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
			_, err := f.ficsitCli.Installations.AddInstallation(f.ficsitCli, install.Path, fallbackProfile)
			if err != nil {
				return errors.Wrap(err, "failed to add installation")
			}
			createdNewInstalls = true
		}
		f.installationMetadata[install.Path] = install
	}

	if createdNewInstalls {
		err := f.ficsitCli.Installations.Save()
		if err != nil {
			return errors.Wrap(err, "failed to save installations")
		}
	}
	return nil
}

func (f *FicsitCLI) initRemoteServerInstallationsMetadata() error {
	for _, installation := range f.GetInstallations() {
		err := f.checkAndAddExistingRemote(installation)
		if err != nil {
			slog.Warn("failed to check and add existing remote", slog.Any("error", err), utils.SlogPath("path", installation))
		}
	}
	return nil
}

func isServerTarget(targetName string) bool {
	return targetName == "WindowsServer" || targetName == "LinuxServer"
}

func (f *FicsitCLI) checkAndAddExistingRemote(path string) error {
	slog.Debug("checking whether installation is remote", utils.SlogPath("path", path))
	installation := f.ficsitCli.Installations.GetInstallation(path)
	if installation == nil {
		return nil
	}
	if _, ok := f.installationMetadata[installation.Path]; ok {
		// Already have metadata for this install
		return nil
	}
	platform, err := installation.GetPlatform(f.ficsitCli)
	if err != nil {
		// Maybe the server is unreachable at the moment
		// We will keep this install for now
		slog.Info("failed to get platform", slog.Any("error", err), utils.SlogPath("path", installation.Path))
	} else if !isServerTarget(platform.TargetName) {
		// Not a server, but a local install, should already have metadata
		return nil
	}
	if err := f.AddRemoteServer(path); err != nil {
		return errors.Wrap(err, "failed to add remote server")
	}
	return nil
}

func (f *FicsitCLI) GetInstallations() []string {
	installations := make([]string, 0, len(f.ficsitCli.Installations.Installations))
	for _, installation := range f.ficsitCli.Installations.Installations {
		// Keep installations that we have metadata for
		if _, ok := f.installationMetadata[installation.Path]; !ok {
			// Keep installations that are remote servers
			// Even if we don't have metadata for them
			platform, err := installation.GetPlatform(f.ficsitCli)
			if err != nil {
				// Maybe the server is unreachable at the moment
				// We will keep this install for now
				slog.Info("failed to get platform", slog.Any("error", err), utils.SlogPath("path", installation.Path))
			} else if !isServerTarget(platform.TargetName) {
				// Not a server, but a local install, should already have metadata
				continue
			}
		}
		installations = append(installations, installation.Path)
	}
	return installations
}

func (f *FicsitCLI) GetInstallationsMetadata() map[string]*common.Installation {
	return f.installationMetadata
}

func (f *FicsitCLI) GetCurrentInstallationMetadata() *common.Installation {
	// This function only exists so common.Installation is generated to typescript
	return f.installationMetadata[f.ficsitCli.Installations.SelectedInstallation]
}

func (f *FicsitCLI) GetInvalidInstalls() []string {
	result := []string{}
	for _, err := range f.installFindErrors {
		var installFindErr common.InstallFindError
		if errors.As(err, &installFindErr) {
			result = append(result, installFindErr.Path)
		}
	}
	return result
}

func (f *FicsitCLI) GetInstallation(path string) *cli.Installation {
	return f.ficsitCli.Installations.GetInstallation(path)
}

func (f *FicsitCLI) SelectInstall(path string) error {
	l := slog.With(slog.String("task", "selectInstall"), utils.SlogPath("path", path))

	if !f.isValidInstall(path) {
		return fmt.Errorf("invalid installation: %s", path)
	}
	if f.ficsitCli.Installations.SelectedInstallation == path {
		return nil
	}
	installation := f.ficsitCli.Installations.GetInstallation(path)
	if installation == nil {
		l.Error("failed to find installation")
		return errors.Errorf("installation %s not found", path)
	}

	f.ficsitCli.Installations.SelectedInstallation = path
	_ = f.ficsitCli.Installations.Save()

	selectedInstallation := f.GetSelectedInstall()

	f.EmitGlobals()

	f.progress = &Progress{
		Item:     "__select_install__",
		Message:  "Validating install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstallation, "__select_install__")

	if installErr != nil {
		l.Error("failed to validate install", slog.Any("error", installErr))
		return installErr
	}
	return nil
}

func (f *FicsitCLI) GetSelectedInstall() *cli.Installation {
	return f.ficsitCli.Installations.GetInstallation(f.ficsitCli.Installations.SelectedInstallation)
}

func (f *FicsitCLI) SetModsEnabled(enabled bool) error {
	selectedInstallation := f.GetSelectedInstall()

	if selectedInstallation == nil {
		slog.Error("no installation selected")
		return errors.New("no installation selected")
	}
	l := slog.With(slog.String("task", "setModsEnabled"), slog.Bool("enabled", enabled), utils.SlogPath("install", selectedInstallation.Path))

	var message string
	if enabled {
		message = "Enabling mods"
	} else {
		message = "Disabling mods"
	}

	selectedInstallation.Vanilla = !enabled
	_ = f.ficsitCli.Installations.Save()

	f.EmitGlobals()

	f.progress = &Progress{
		Item:     "__toggle_mods__",
		Message:  message,
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstallation, "__toggle_mods__")

	if installErr != nil {
		l.Error("failed to validate install", slog.Any("error", installErr))
		return installErr
	}

	return nil
}

func (f *FicsitCLI) GetModsEnabled() bool {
	selectedInstallation := f.GetSelectedInstall()
	return selectedInstallation == nil || !selectedInstallation.Vanilla
}

func (f *FicsitCLI) GetSelectedInstallProfileMods() map[string]cli.ProfileMod {
	selectedInstallation := f.GetSelectedInstall()
	if selectedInstallation == nil {
		return make(map[string]cli.ProfileMod)
	}
	profile := f.GetProfile(selectedInstallation.Profile)
	return profile.Mods
}

func (f *FicsitCLI) GetSelectedInstallLockfileMods() (map[string]resolver.LockedMod, error) {
	selectedInstallation := f.GetSelectedInstall()
	if selectedInstallation == nil {
		return make(map[string]resolver.LockedMod), nil
	}
	lockfile, err := selectedInstallation.LockFile(f.ficsitCli)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	if lockfile == nil {
		return make(map[string]resolver.LockedMod), nil
	}
	return lockfile.Mods, nil
}

func (f *FicsitCLI) GetSelectedInstallLockfile() (*resolver.LockFile, error) {
	selectedInstallation := f.GetSelectedInstall()
	if selectedInstallation == nil {
		return nil, nil
	}
	lockfile, err := selectedInstallation.LockFile(f.ficsitCli)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return lockfile, nil
}

func (f *FicsitCLI) LaunchGame() {
	selectedInstallation := f.GetSelectedInstall()
	if selectedInstallation == nil {
		slog.Error("no installation selected")
		return
	}
	metadata := f.installationMetadata[selectedInstallation.Path]
	if metadata == nil {
		slog.Error("no metadata for installation")
		return
	}
	cmd := exec.Command(metadata.LaunchPath[0], metadata.LaunchPath[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("failed to launch game", slog.Any("error", err), slog.String("cmd", cmd.String()), slog.String("output", string(out)))
		return
	}
}
