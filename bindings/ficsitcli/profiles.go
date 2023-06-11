package ficsitcli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/settings"
)

func (f *FicsitCLI) SetProfile(profile string) error {
	l := log.With().Str("task", "setProfile").Str("profile", profile).Logger()

	if f.selectedInstallation == nil {
		l.Error().Str("profile", profile).Msg("No installation selected")
		return errors.New("No installation selected")
	}
	if f.selectedInstallation.Installation.Profile == profile {
		return nil
	}
	err := f.selectedInstallation.Installation.SetProfile(f.ficsitCli, profile)
	if err != nil {
		l.Error().Err(err).Str("profile", profile).Msg("Failed to set profile")
		return errors.Wrap(err, "Failed to set profile")
	}

	f.progress = &Progress{
		Item:     "__select_profile__",
		Message:  "Validating install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(f.selectedInstallation, "__select_profile__")

	if installErr != nil {
		l.Error().Err(installErr).Str("profile", profile).Msg("Failed to validate install")
		return errors.Wrap(installErr, "Failed to validate install")
	}

	settings.Settings.SelectedProfile[f.selectedInstallation.Info.Path] = profile
	_ = settings.SaveSettings()
	return nil
}

func (f *FicsitCLI) GetSelectedProfile() *string {
	if f.selectedInstallation == nil {
		return nil
	}
	return &f.selectedInstallation.Installation.Profile
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
	l := log.With().Str("task", "addProfile").Str("profile", name).Logger()

	_, err := f.ficsitCli.Profiles.AddProfile(name)
	if err != nil {
		l.Error().Err(err).Str("name", name).Msg("Failed to add profile")
		return errors.Wrapf(err, "Failed to add profile: %s", name)
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) RenameProfile(oldName string, newName string) error {
	l := log.With().Str("task", "renameProfile").Str("oldName", oldName).Str("newName", newName).Logger()

	err := f.ficsitCli.Profiles.RenameProfile(f.ficsitCli, oldName, newName)
	if err != nil {
		l.Error().Err(err).Str("oldName", oldName).Str("newName", newName).Msg("Failed to rename profile")
		return errors.Wrapf(err, "Failed to rename profile: %s -> %s", oldName, newName)
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) DeleteProfile(name string) error {
	l := log.With().Str("task", "deleteProfile").Str("profile", name).Logger()

	err := f.ficsitCli.Profiles.DeleteProfile(name)
	if err != nil {
		l.Error().Err(err).Str("name", name).Msg("Failed to delete profile")
		return errors.Wrapf(err, "Failed to delete profile: %s", name)
	}

	_ = f.ficsitCli.Profiles.Save()

	for _, install := range f.installations {
		if install.Installation.Profile == name {
			install.Installation.Profile = cli.DefaultProfileName
		}
	}

	return nil
}

type ExportedProfile struct {
	Profile  *cli.Profile             `json:"profile"`
	LockFile *cli.LockFile            `json:"lockfile"`
	Metadata *ExportedProfileMetadata `json:"metadata"`
}

type ExportedProfileMetadata struct {
	GameVersion int `json:"gameVersion"`
}

func (f *FicsitCLI) MakeCurrentExportedProfile() (*ExportedProfile, error) {
	l := log.With().Str("task", "makeCurrentExportedProfile").Logger()

	selectedInstall := f.GetSelectedInstall()
	if selectedInstall == nil {
		l.Error().Msg("No installation selected")
		return nil, errors.New("No installation selected")
	}

	profileName := f.GetSelectedProfile()
	if profileName == nil {
		l.Error().Msg("No profile selected")
		return nil, errors.New("No profile selected")
	}

	profile := f.GetProfile(*profileName)
	if profile == nil {
		l.Error().Str("profile", *profileName).Msg("Profile not found")
		return nil, errors.New("Profile not found")
	}
	lockfile, err := selectedInstall.Installation.LockFile(f.ficsitCli)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get lockfile")
		return nil, errors.Wrap(err, "Failed to get lockfile")
	}
	metadata := &ExportedProfileMetadata{
		GameVersion: selectedInstall.Info.Version,
	}

	return &ExportedProfile{
		Profile:  profile,
		LockFile: lockfile,
		Metadata: metadata,
	}, nil
}

func (f *FicsitCLI) ExportCurrentProfile() error {
	l := log.With().Str("task", "exportCurrentProfile").Logger()

	exportedProfile, err := f.MakeCurrentExportedProfile()
	if err != nil {
		l.Error().Err(err).Msg("Failed to make exported profile")
		return errors.Wrapf(err, "Failed to export profile")
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
		l.Error().Err(err).Msg("Failed to open save dialog")
		return errors.Wrapf(err, "Failed to export profile: %s", exportedProfile.Profile.Name)
	}

	exportedProfileJSON, err := json.MarshalIndent(exportedProfile, "", "  ")
	if err != nil {
		l.Error().Err(err).Msg("Failed to marshal exported profile")
		return errors.Wrapf(err, "Failed to export profile: %s", exportedProfile.Profile.Name)
	}
	err = os.WriteFile(filename, exportedProfileJSON, 0o755)
	if err != nil {
		l.Error().Err(err).Msg("Failed to write exported profile")
		return errors.Wrapf(err, "Failed to export profile: %s", exportedProfile.Profile.Name)
	}

	return nil
}

func (f *FicsitCLI) ReadExportedProfileMetadata(file string) (*ExportedProfileMetadata, error) {
	l := log.With().Str("task", "readExportedProfileMetadata").Str("file", file).Logger()

	fileBytes, err := os.ReadFile(file)
	if err != nil {
		l.Error().Err(err).Str("file", file).Msg("Failed to read exported profile")
		return nil, errors.Wrap(err, "Failed to read exported profile")
	}

	var exportedProfile ExportedProfile
	err = json.Unmarshal(fileBytes, &exportedProfile)
	if err != nil {
		l.Error().Err(err).Str("file", file).Msg("Failed to unmarshal exported profile")
		return nil, errors.Wrap(err, "Failed to read exported profile")
	}

	return exportedProfile.Metadata, nil
}

func (f *FicsitCLI) ImportProfile(name string, file string) error {
	l := log.With().Str("task", "importProfile").Str("name", name).Str("file", file).Logger()

	selectedInstall := f.GetSelectedInstall()
	if selectedInstall == nil {
		l.Error().Msg("No installation selected")
		return errors.New("No installation selected")
	}

	profileData, err := os.ReadFile(file)
	if err != nil {
		l.Error().Err(err).Str("file", file).Msg("Failed to read exported profile")
		return errors.Wrap(err, "Failed to read profile file")
	}

	var exportedProfile ExportedProfile
	err = json.Unmarshal(profileData, &exportedProfile)
	if err != nil {
		l.Error().Err(err).Str("file", file).Msg("Failed to unmarshal exported profile")
		return errors.Wrap(err, "Failed to read profile file")
	}

	profile, err := f.ficsitCli.Profiles.AddProfile(name)
	if err != nil {
		l.Error().Err(err).Str("name", name).Msg("Failed to add profile")
		return errors.Wrap(err, "Failed to add imported profile")
	}

	profile.Mods = exportedProfile.Profile.Mods

	_ = selectedInstall.Installation.SetProfile(f.ficsitCli, name)

	err = selectedInstall.Installation.WriteLockFile(f.ficsitCli, *exportedProfile.LockFile)
	if err != nil {
		_ = f.ficsitCli.Profiles.DeleteProfile(name)
		l.Error().Err(err).Str("name", name).Msg("Failed to write lockfile")
		return errors.Wrap(err, "Failed to write profile")
	}

	f.progress = &Progress{
		Item:     "__import_profile__",
		Message:  "Validating install",
		Progress: -1,
	}

	f.setProgress(f.progress)

	defer f.setProgress(nil)

	installErr := f.validateInstall(selectedInstall, "__import_profile__")

	if installErr != nil {
		_ = f.ficsitCli.Profiles.DeleteProfile(name)
		l.Error().Err(installErr).Str("profile", name).Msg("Failed to validate install")
		return errors.Wrap(installErr, "Failed to validate install")
	}

	_ = f.ficsitCli.Profiles.Save()

	return nil
}
