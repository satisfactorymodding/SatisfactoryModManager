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

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

func (f *FicsitCLI) SetProfile(profile string) error {
	l := slog.With(slog.String("task", "setProfile"), slog.String("profile", profile))

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
	_ = f.ficsitCli.Installations.Save()

	f.EmitGlobals()

	f.progress = &Progress{
		Item:     "__select_profile__",
		Message:  "Validating install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstallation, "__select_profile__")

	if installErr != nil {
		l.Error("failed to validate installation", slog.Any("error", installErr))
		return installErr
	}

	return nil
}

func (f *FicsitCLI) GetSelectedProfile() *string {
	selectedInstallation := f.GetSelectedInstall()
	if selectedInstallation == nil {
		return nil
	}
	return &selectedInstallation.Profile
}

func (f *FicsitCLI) GetProfiles() []string {
	profileNames := make([]string, 0, len(f.ficsitCli.Profiles.Profiles))
	for k := range f.ficsitCli.Profiles.Profiles {
		profileNames = append(profileNames, k)
	}
	sort.Strings(profileNames)
	return profileNames
}

func (f *FicsitCLI) GetProfile(profile string) *cli.Profile {
	return f.ficsitCli.Profiles.GetProfile(profile)
}

func (f *FicsitCLI) AddProfile(name string) error {
	l := slog.With(slog.String("task", "addProfile"), slog.String("profile", name))

	_, err := f.ficsitCli.Profiles.AddProfile(name)
	if err != nil {
		l.Error("failed to add profile", slog.Any("error", err))
		return fmt.Errorf("failed to add profile: %s: %w", name, err)
	}

	_ = f.ficsitCli.Profiles.Save()

	f.EmitGlobals()

	return nil
}

func (f *FicsitCLI) RenameProfile(oldName string, newName string) error {
	l := slog.With(slog.String("task", "renameProfile"), slog.String("oldName", oldName), slog.String("newName", newName))

	err := f.ficsitCli.Profiles.RenameProfile(f.ficsitCli, oldName, newName)
	if err != nil {
		l.Error("failed to rename profile", slog.Any("error", err))
		return fmt.Errorf("failed to rename profile: %s -> %s: %w", oldName, newName, err)
	}

	_ = f.ficsitCli.Profiles.Save()
	_ = f.ficsitCli.Installations.Save() // Installs using the old name will be updated

	f.EmitGlobals()

	return nil
}

func (f *FicsitCLI) DeleteProfile(name string) error {
	l := slog.With(slog.String("task", "deleteProfile"), slog.String("profile", name))

	err := f.ficsitCli.Profiles.DeleteProfile(name)
	if err != nil {
		l.Error("failed to delete profile", slog.Any("error", err))
		return fmt.Errorf("failed to delete profile: %s: %w", name, err)
	}

	_ = f.ficsitCli.Profiles.Save()

	for _, install := range f.ficsitCli.Installations.Installations {
		if install.Profile == name {
			install.Profile = cli.DefaultProfileName
		}
	}

	_ = f.ficsitCli.Installations.Save()

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

func (f *FicsitCLI) MakeCurrentExportedProfile() (*ExportedProfile, error) {
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

	installMetadata := f.installationMetadata[selectedInstallation.Path]
	var gameVersion int
	if installMetadata != nil {
		gameVersion = installMetadata.Version
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

func (f *FicsitCLI) ExportCurrentProfile() error {
	l := slog.With(slog.String("task", "exportCurrentProfile"))

	exportedProfile, err := f.MakeCurrentExportedProfile()
	if err != nil {
		l.Error("failed to make exported profile", slog.Any("error", err))
		return fmt.Errorf("failed to export profile: %w", err)
	}

	defaultFileName := fmt.Sprintf("%s-%s.smmprofile", exportedProfile.Profile.Name, time.Now().UTC().Format("2006-01-02-15-04-05"))
	filename, err := wailsRuntime.SaveFileDialog(f.ctx, wailsRuntime.SaveDialogOptions{
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

func (f *FicsitCLI) ReadExportedProfileMetadata(file string) (*ExportedProfileMetadata, error) {
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

func (f *FicsitCLI) ImportProfile(name string, file string) error {
	l := slog.With(slog.String("task", "importProfile"), slog.String("name", name), slog.String("file", file))

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

	f.progress = &Progress{
		Item:     "__import_profile__",
		Message:  "Validating install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstallation, "__import_profile__")

	if installErr != nil {
		_ = f.ficsitCli.Profiles.DeleteProfile(name)
		l.Error("failed to validate installation", slog.Any("error", installErr))
		return installErr
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}
