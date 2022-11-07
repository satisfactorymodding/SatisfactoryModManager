package ficsitcli_bindings

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/SatisfactoryModManager/settings"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (f *FicsitCLI) SetProfile(profile string) error {
	if f.selectedInstallation == nil {
		log.Error().Str("profile", profile).Msg("No installation selected")
		return errors.New("No installation selected")
	}
	err := f.selectedInstallation.Installation.SetProfile(f.ficsitCli, profile)
	if err != nil {
		log.Error().Err(err).Str("profile", profile).Msg("Failed to set profile")
		return errors.Wrap(err, "Failed to set profile")
	}
	f.emitModsChange()
	settings.Settings.SelectedProfile[f.selectedInstallation.Info.Path] = profile
	settings.SaveSettings()
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
	_, err := f.ficsitCli.Profiles.AddProfile(name)

	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to add profile")
		return errors.Wrapf(err, "Failed to add profile: %s", name)
	}

	f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) RenameProfile(oldName string, newName string) error {
	err := f.ficsitCli.Profiles.RenameProfile(f.ficsitCli, oldName, newName)

	if err != nil {
		log.Error().Err(err).Str("oldName", oldName).Str("newName", newName).Msg("Failed to rename profile")
		return errors.Wrapf(err, "Failed to rename profile: %s -> %s", oldName, newName)
	}

	f.ficsitCli.Profiles.Save()

	return nil
}

func (f *FicsitCLI) DeleteProfile(name string) error {
	err := f.ficsitCli.Profiles.DeleteProfile(name)

	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to delete profile")
		return errors.Wrapf(err, "Failed to delete profile: %s", name)
	}

	f.ficsitCli.Profiles.Save()

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
	selectedInstall := f.GetSelectedInstall()
	if selectedInstall == nil {
		log.Error().Msg("No installation selected")
		return nil, errors.New("No installation selected")
	}

	profileName := f.GetSelectedProfile()
	if profileName == nil {
		log.Error().Msg("No profile selected")
		return nil, errors.New("No profile selected")
	}

	profile := f.GetProfile(*profileName)
	if profile == nil {
		log.Error().Str("profile", *profileName).Msg("Profile not found")
		return nil, errors.New("Profile not found")
	}
	lockfile, err := selectedInstall.Installation.LockFile(f.ficsitCli)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get lockfile")
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
	exportedProfile, err := f.MakeCurrentExportedProfile()

	if err != nil {
		log.Error().Err(err).Msg("Failed to make exported profile")
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
		log.Error().Err(err).Msg("Failed to open save dialog")
		return errors.Wrapf(err, "Failed to export profile: %s", exportedProfile.Profile.Name)
	}

	exportedProfileJson, err := json.MarshalIndent(exportedProfile, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal exported profile")
		return errors.Wrapf(err, "Failed to export profile: %s", exportedProfile.Profile.Name)
	}
	err = os.WriteFile(filename, exportedProfileJson, 0755)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write exported profile")
		return errors.Wrapf(err, "Failed to export profile: %s", exportedProfile.Profile.Name)
	}

	return nil
}

func (f *FicsitCLI) ReadExportedProfileMetadata(file string) (*ExportedProfileMetadata, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		log.Error().Err(err).Str("file", file).Msg("Failed to read exported profile")
		return nil, errors.Wrap(err, "Failed to read exported profile")
	}

	var exportedProfile ExportedProfile
	err = json.Unmarshal(fileBytes, &exportedProfile)
	if err != nil {
		log.Error().Err(err).Str("file", file).Msg("Failed to unmarshal exported profile")
		return nil, errors.Wrap(err, "Failed to read exported profile")
	}

	return exportedProfile.Metadata, nil
}

func (f *FicsitCLI) ImportProfile(name string, file string) error {
	selectedInstall := f.GetSelectedInstall()
	if selectedInstall == nil {
		log.Error().Msg("No installation selected")
		return errors.New("No installation selected")
	}

	profileData, err := os.ReadFile(file)
	if err != nil {
		log.Error().Err(err).Str("file", file).Msg("Failed to read exported profile")
		return errors.Wrap(err, "Failed to read profile file")
	}

	var exportedProfile ExportedProfile
	err = json.Unmarshal(profileData, &exportedProfile)
	if err != nil {
		log.Error().Err(err).Str("file", file).Msg("Failed to unmarshal exported profile")
		return errors.Wrap(err, "Failed to read profile file")
	}

	profile, err := f.ficsitCli.Profiles.AddProfile(name)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to add profile")
		return errors.Wrap(err, "Failed to add imported profile")
	}

	profile.Mods = exportedProfile.Profile.Mods

	selectedInstall.Installation.SetProfile(f.ficsitCli, name)

	err = selectedInstall.Installation.WriteLockFile(f.ficsitCli, *exportedProfile.LockFile)
	if err != nil {
		f.ficsitCli.Profiles.DeleteProfile(name)
		log.Error().Err(err).Str("name", name).Msg("Failed to write lockfile")
		return errors.Wrap(err, "Failed to write profile")
	}

	f.ficsitCli.Profiles.Save()

	return nil
}
