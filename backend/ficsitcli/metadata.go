package ficsitcli

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/satisfactorymodding/ficsit-cli/cli"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
)

func (f *ficsitCLI) initLocalInstallationsMetadata() error {
	installs, findErrors := installfinders.FindInstallations()

	f.installFindErrors = findErrors

	for _, findErr := range findErrors {
		slog.Info("failed to find installations", slog.Any("error", findErr))
	}

	fallbackProfile := f.GetFallbackProfile()

	createdNewInstalls := false
	for _, install := range installs {
		ficsitCliInstall := f.ficsitCli.Installations.GetInstallation(install.Path)
		if ficsitCliInstall == nil {
			_, err := f.ficsitCli.Installations.AddInstallation(f.ficsitCli, install.Path, fallbackProfile)
			if err != nil {
				return fmt.Errorf("failed to add installation: %w", err)
			}
			createdNewInstalls = true
		} else {
			_, profileExists := f.ficsitCli.Profiles.Profiles[ficsitCliInstall.Profile]
			if !profileExists {
				ficsitCliInstall.Profile = fallbackProfile
				createdNewInstalls = true
			}
		}
		f.installationMetadata.Store(install.Path, installationMetadata{
			State: InstallStateValid,
			Info:  install,
		})
	}

	if createdNewInstalls {
		err := f.ficsitCli.Installations.Save()
		if err != nil {
			return fmt.Errorf("failed to save installations: %w", err)
		}
	}
	return nil
}

func (f *ficsitCLI) initRemoteServerInstallationsMetadata() {
	installationsToCheck := make([]*cli.Installation, 0, len(f.ficsitCli.Installations.Installations))
	for _, installation := range f.ficsitCli.Installations.Installations {
		if meta, ok := f.installationMetadata.Load(installation.Path); ok {
			if meta.State != InstallStateUnknown {
				// Already have metadata for this install
				continue
			}
		}
		installationsToCheck = append(installationsToCheck, installation)
	}

	if len(installationsToCheck) == 0 {
		return
	}

	// Mark all installations as loading, since this may take a while

	for _, installation := range installationsToCheck {
		f.installationMetadata.Store(installation.Path, installationMetadata{
			State: InstallStateLoading,
		})
	}

	f.EmitGlobals()

	var wg sync.WaitGroup
	for _, installation := range installationsToCheck {
		wg.Add(1)
		go func(installation *cli.Installation) {
			defer wg.Done()
			f.fetchRemoteInstallationMetadata(installation)
		}(installation)
	}
	wg.Wait()

	meta, ok := f.installationMetadata.Load(f.ficsitCli.Installations.SelectedInstallation)
	if !ok || meta.State == InstallStateInvalid {
		f.ensureSelectedInstallationIsValid()
	}

	f.EmitGlobals()
}

func (f *ficsitCLI) fetchRemoteInstallationMetadata(installation *cli.Installation) {
	defer f.EmitGlobals()
	meta, err := f.getRemoteServerMetadata(installation)
	if err != nil {
		if errors.Is(err, ErrInstallNotServer) {
			// If this installation is not a server, it is invalid
			f.installationMetadata.Store(installation.Path, installationMetadata{
				State: InstallStateInvalid,
			})
			return
		}
		// If we failed to get metadata, we will keep this install for now
		f.installationMetadata.Store(installation.Path, installationMetadata{
			State: InstallStateUnknown,
			Info:  nil,
		})
		slog.Warn("failed to get remote server metadata", slog.Any("error", err), slog.String("path", installation.Path))
		return
	}

	f.installationMetadata.Store(installation.Path, installationMetadata{
		State: InstallStateValid,
		Info:  meta,
	})
}

func (f *ficsitCLI) FetchRemoteServerMetadata(path string) error {
	installation := f.ficsitCli.Installations.GetInstallation(path)
	if installation == nil {
		return fmt.Errorf("installation not found")
	}
	if meta, ok := f.installationMetadata.Load(path); ok && meta.State != InstallStateUnknown {
		return nil
	}
	f.installationMetadata.Store(path, installationMetadata{
		State: InstallStateLoading,
	})
	f.EmitGlobals()
	f.fetchRemoteInstallationMetadata(installation)
	return nil
}

var ErrInstallNotServer = fmt.Errorf("installation is not a server")

func (f *ficsitCLI) getRemoteServerMetadata(installation *cli.Installation) (*common.Installation, error) {
	gameVersion, err := installation.GetGameVersion(f.ficsitCli)
	if err != nil {
		return nil, fmt.Errorf("failed to get game version: %w", err)
	}

	platform, err := installation.GetPlatform(f.ficsitCli)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform: %w", err)
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
		return nil, ErrInstallNotServer
	}

	branch := common.BranchStable // TODO: Do we have a way to detect this for remote installs?

	remoteName := settings.Settings.RemoteNames[remoteKey(installation.Path)]

	if remoteName == "" {
		remoteName = f.GetNextRemoteLauncherName()
	}

	return &common.Installation{
		Path:      installation.Path,
		Type:      installType,
		Location:  common.LocationTypeRemote,
		Branch:    branch,
		Version:   gameVersion,
		Launcher:  remoteName,
		SavedPath: path.Join(installation.BasePath(), "FactoryGame", "Saved"),
	}, nil
}

func remoteKey(path string) string {
	hash := sha256.Sum256([]byte(path))
	return hex.EncodeToString(hash[:])
}

func (f *ficsitCLI) GetNextRemoteLauncherName() string {
	existingNumbers := make(map[int]bool)
	for _, install := range f.GetRemoteInstallations() {
		metadata, ok := f.installationMetadata.Load(install)
		if ok && metadata.Info != nil {
			if strings.HasPrefix(metadata.Info.Launcher, "Server ") {
				num, err := strconv.Atoi(strings.TrimPrefix(metadata.Info.Launcher, "Server "))
				if err == nil {
					existingNumbers[num] = true
				}
			}
		}
	}
	for i := 1; ; i++ {
		if !existingNumbers[i] {
			return "Server " + strconv.Itoa(i)
		}
	}
}
