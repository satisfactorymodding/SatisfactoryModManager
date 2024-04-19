package ficsitcli

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"time"

	"github.com/satisfactorymodding/ficsit-cli/cli"
	resolver "github.com/satisfactorymodding/ficsit-resolver"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	appCommon "github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

func (f *ficsitCLI) SetProfile(profile string) error {
	return f.action(ActionSelectProfile, newSimpleItem(profile), func(l *slog.Logger, taskChannel chan<- taskUpdate) error {
		selectedInstallation := f.GetSelectedInstall()

		if selectedInstallation == nil {
			l.Error("no installation selected")
			return fmt.Errorf("no installation selected")
		}
		if selectedInstallation.Profile == profile {
			return nil
		}

		err := selectedInstallation.SetProfile(f.ficsitCli, profile)
		if err != nil {
			l.Error("failed to set profile", slog.Any("error", err))
			return fmt.Errorf("failed to set profile: %w", err)
		}

		err = f.ficsitCli.Installations.Save()
		if err != nil {
			l.Error("failed to save installations", slog.Any("error", err))
		}

		f.EmitGlobals()

		installErr := f.validateInstall(selectedInstallation, taskChannel)

		if installErr != nil {
			l.Error("failed to validate installation", slog.Any("error", installErr))
			return installErr
		}

		return nil
	})
}

func (f *ficsitCLI) GetSelectedProfile() *string {
	selectedInstallation := f.GetSelectedInstall()
	if selectedInstallation == nil {
		return nil
	}
	return &selectedInstallation.Profile
}

func (f *ficsitCLI) GetProfiles() []string {
	profileNames := make([]string, 0, len(f.ficsitCli.Profiles.Profiles))
	for k := range f.ficsitCli.Profiles.Profiles {
		profileNames = append(profileNames, k)
	}
	sort.Strings(profileNames)
	return profileNames
}

func (f *ficsitCLI) GetProfile(profile string) *cli.Profile {
	return f.ficsitCli.Profiles.GetProfile(profile)
}

func (f *ficsitCLI) GetFallbackProfile() string {
	fallbackProfile := cli.DefaultProfileName
	if f.ficsitCli.Profiles.GetProfile(fallbackProfile) == nil {
		// Pick first profile found
		for name := range f.ficsitCli.Profiles.Profiles {
			fallbackProfile = name
			break
		}
	}
	return fallbackProfile
}

func (f *ficsitCLI) GetFallbackProfileExcept(profile string) string {
	fallbackProfile := cli.DefaultProfileName
	if f.ficsitCli.Profiles.GetProfile(fallbackProfile) == nil {
		// Pick first profile found, that is not excluded
		for name := range f.ficsitCli.Profiles.Profiles {
			if name == profile {
				continue
			}
			fallbackProfile = name
			break
		}
	}
	return fallbackProfile
}

func (f *ficsitCLI) AddProfile(name string) error {
	l := slog.With(slog.String("task", "addProfile"), slog.String("profile", name))

	_, err := f.ficsitCli.Profiles.AddProfile(name)
	if err != nil {
		l.Error("failed to add profile", slog.Any("error", err))
		return fmt.Errorf("failed to add profile: %s: %w", name, err)
	}

	err = f.ficsitCli.Profiles.Save()
	if err != nil {
		l.Error("failed to save profile", slog.Any("error", err))
	}

	f.EmitGlobals()

	return nil
}

func (f *ficsitCLI) RenameProfile(oldName string, newName string) error {
	l := slog.With(slog.String("task", "renameProfile"), slog.String("oldName", oldName), slog.String("newName", newName))

	err := f.ficsitCli.Profiles.RenameProfile(f.ficsitCli, oldName, newName)
	if err != nil {
		l.Error("failed to rename profile", slog.Any("error", err))
		return fmt.Errorf("failed to rename profile: %s -> %s: %w", oldName, newName, err)
	}

	err = f.ficsitCli.Profiles.Save()
	if err != nil {
		l.Error("failed to save profile", slog.Any("error", err))
	}

	// Installs using the old name will be updated
	err = f.ficsitCli.Installations.Save()
	if err != nil {
		l.Error("failed to save installations", slog.Any("error", err))
	}

	f.EmitGlobals()

	return nil
}

func (f *ficsitCLI) DeleteProfile(name string) error {
	l := slog.With(slog.String("task", "deleteProfile"), slog.String("profile", name))

	// ficsit-cli always sets installs that use the deleted profile to Default, which might not exist
	fallbackProfile := f.GetFallbackProfileExcept(name)
	for _, installation := range f.ficsitCli.Installations.Installations {
		if installation.Profile == name {
			_ = installation.SetProfile(f.ficsitCli, fallbackProfile)
		}
	}

	err := f.ficsitCli.Profiles.DeleteProfile(name)
	if err != nil {
		l.Error("failed to delete profile", slog.Any("error", err))
		return fmt.Errorf("failed to delete profile: %s: %w", name, err)
	}

	err = f.ficsitCli.Profiles.Save()
	if err != nil {
		l.Error("failed to save profile", slog.Any("error", err))
	}

	// Installs using the profile will be updated
	err = f.ficsitCli.Installations.Save()
	if err != nil {
		l.Error("failed to save installations", slog.Any("error", err))
	}

	f.EmitGlobals()

	return nil
}

type ExportedProfile struct {
	Profile  cli.Profile              `json:"profile"`
	LockFile resolver.LockFile        `json:"lockfile"`
	Metadata *ExportedProfileMetadata `json:"metadata"`
}

type ExportedProfileMetadata struct {
	GameVersion int `json:"gameVersion"`
}

func (f *ficsitCLI) MakeCurrentExportedProfile() (*ExportedProfile, error) {
	l := slog.With(slog.String("task", "makeCurrentExportedProfile"))

	selectedInstallation := f.GetSelectedInstall()

	if selectedInstallation == nil {
		l.Error("no installation selected")
		return nil, fmt.Errorf("no installation selected")
	}

	profileName := f.GetSelectedProfile()
	if profileName == nil {
		l.Error("no profile selected")
		return nil, fmt.Errorf("no profile selected")
	}

	profile := f.GetProfile(*profileName)
	if profile == nil {
		l.Error("profile not found", slog.String("profile", *profileName))
		return nil, fmt.Errorf("profile not found")
	}
	lockfile, err := selectedInstallation.LockFile(f.ficsitCli)
	if err != nil {
		l.Error("failed to get lockfile", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get lockfile: %w", err)
	}

	installMetadata, ok := f.installationMetadata.Load(selectedInstallation.Path)
	var gameVersion int
	if ok && installMetadata.Info != nil {
		gameVersion = installMetadata.Info.Version
	}
	metadata := &ExportedProfileMetadata{
		GameVersion: gameVersion,
	}

	if lockfile == nil {
		lockfile = resolver.NewLockfile()
	}

	return &ExportedProfile{
		Profile:  *profile,
		LockFile: *lockfile,
		Metadata: metadata,
	}, nil
}

func (f *ficsitCLI) ExportCurrentProfile() error {
	l := slog.With(slog.String("task", "exportCurrentProfile"))

	exportedProfile, err := f.MakeCurrentExportedProfile()
	if err != nil {
		l.Error("failed to make exported profile", slog.Any("error", err))
		return fmt.Errorf("failed to export profile: %w", err)
	}

	defaultFileName := fmt.Sprintf("%s-%s.smmprofile", exportedProfile.Profile.Name, time.Now().UTC().Format("2006-01-02-15-04-05"))
	filename, err := wailsRuntime.SaveFileDialog(appCommon.AppContext, wailsRuntime.SaveDialogOptions{
		DefaultFilename: defaultFileName,
		Filters: []wailsRuntime.FileFilter{
			{
				Pattern:     "*.smmprofile",
				DisplayName: "SMM Profile (*.smmprofile)",
			},
		},
	})
	if err != nil {
		l.Error("failed to open save dialog", slog.Any("error", err))
		return fmt.Errorf("failed to open save dialog: %w", err)
	}
	if filename == "" {
		// User cancelled
		return nil
	}

	exportedProfileJSON, err := utils.JSONMarshal(exportedProfile, 2)
	if err != nil {
		l.Error("failed to marshal exported profile", slog.Any("error", err))
		return fmt.Errorf("failed to marshal exported profile: %w", err)
	}
	err = os.WriteFile(filename, exportedProfileJSON, 0o755)
	if err != nil {
		l.Error("failed to write exported profile", slog.Any("error", err))
		return fmt.Errorf("failed to write exported profile: %w", err)
	}

	return nil
}

func (f *ficsitCLI) ReadExportedProfileMetadata(file string) (*ExportedProfileMetadata, error) {
	l := slog.With(slog.String("task", "readExportedProfileMetadata"), slog.String("file", file))

	fileBytes, err := os.ReadFile(file)
	if err != nil {
		l.Error("failed to read exported profile", slog.Any("error", err))
		return nil, fmt.Errorf("failed to read exported profile: %w", err)
	}

	var exportedProfile ExportedProfile
	err = json.Unmarshal(fileBytes, &exportedProfile)
	if err != nil {
		l.Error("failed to unmarshal exported profile", slog.Any("error", err))
		return nil, fmt.Errorf("failed to parse exported profile: %w", err)
	}

	return exportedProfile.Metadata, nil
}

func (f *ficsitCLI) ImportProfile(name string, file string) error {
	return f.action(ActionImportProfile, newSimpleItem(name), func(l *slog.Logger, taskChannel chan<- taskUpdate) error {
		l = l.With(slog.String("file", file))

		selectedInstallation := f.GetSelectedInstall()

		if selectedInstallation == nil {
			l.Error("no installation selected")
			return fmt.Errorf("no installation selected")
		}

		profileData, err := os.ReadFile(file)
		if err != nil {
			l.Error("failed to read exported profile", slog.Any("error", err))
			return fmt.Errorf("failed to read profile file: %w", err)
		}

		var exportedProfile ExportedProfile
		err = json.Unmarshal(profileData, &exportedProfile)
		if err != nil {
			l.Error("failed to unmarshal exported profile", slog.Any("error", err))
			return fmt.Errorf("failed to read profile file: %w", err)
		}

		profile, err := f.ficsitCli.Profiles.AddProfile(name)
		if err != nil {
			l.Error("failed to add profile", slog.Any("error", err))
			return fmt.Errorf("failed to add imported profile: %w", err)
		}

		profile.Mods = exportedProfile.Profile.Mods

		currentProfile := selectedInstallation.Profile

		_ = selectedInstallation.SetProfile(f.ficsitCli, name)

		err = selectedInstallation.WriteLockFile(f.ficsitCli, &exportedProfile.LockFile)
		if err != nil {
			_ = selectedInstallation.SetProfile(f.ficsitCli, currentProfile)
			_ = f.ficsitCli.Profiles.DeleteProfile(name)
			l.Error("failed to write lockfile", slog.Any("error", err))
			return fmt.Errorf("failed to write profile: %w", err)
		}

		f.EmitGlobals()

		installErr := f.validateInstall(selectedInstallation, taskChannel)

		if installErr != nil {
			_ = f.ficsitCli.Profiles.DeleteProfile(name)
			l.Error("failed to validate installation", slog.Any("error", installErr))
			return installErr
		}

		err = f.ficsitCli.Profiles.Save()
		if err != nil {
			l.Error("failed to save profile", slog.Any("error", err))
		}

		return nil
	})
}
