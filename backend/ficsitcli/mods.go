package ficsitcli

import (
	"fmt"
	"log/slog"
)

func (f *ficsitCLI) InstallMod(mod string) error {
	return f.action(ActionInstall, newSimpleItem(mod), func(l *slog.Logger, taskUpdates chan<- taskUpdate) error {
		selectedInstallation := f.GetSelectedInstall()

		if selectedInstallation == nil {
			return fmt.Errorf("no installation selected")
		}

		l = l.With(
			slog.String("install", selectedInstallation.Path),
			slog.String("profile", selectedInstallation.Profile),
		)

		profileName := selectedInstallation.Profile
		profile := f.GetProfile(profileName)

		profileErr := profile.AddMod(mod, ">=0.0.0")
		if profileErr != nil {
			l.Error("failed to add mod", slog.Any("error", profileErr))
			return fmt.Errorf("failed to add mod: %s@latest: %w", mod, profileErr)
		}

		err := f.ficsitCli.Profiles.Save()
		if err != nil {
			l.Error("failed to save profile", slog.Any("error", err))
		}

		installErr := f.validateInstall(selectedInstallation, taskUpdates)

		if installErr != nil {
			l.Error("failed to install", slog.Any("error", installErr))
			return installErr
		}

		return nil
	})
}

func (f *ficsitCLI) InstallModVersion(mod string, version string) error {
	return f.action(ActionInstall, newItem(mod, version), func(l *slog.Logger, taskUpdates chan<- taskUpdate) error {
		selectedInstallation := f.GetSelectedInstall()

		if selectedInstallation == nil {
			return fmt.Errorf("no installation selected")
		}

		l = l.With(
			slog.String("install", selectedInstallation.Path),
			slog.String("profile", selectedInstallation.Profile),
		)

		profile := f.GetProfile(selectedInstallation.Profile)

		profileErr := profile.AddMod(mod, version)
		if profileErr != nil {
			l.Error("failed to add mod", slog.Any("error", profileErr))
			return fmt.Errorf("failed to add mod: %s@%s: %w", mod, version, profileErr)
		}

		err := f.ficsitCli.Profiles.Save()
		if err != nil {
			l.Error("failed to save profile", slog.Any("error", err))
		}

		installErr := f.validateInstall(selectedInstallation, taskUpdates)

		if installErr != nil {
			l.Error("failed to install", slog.Any("error", installErr))
			return installErr
		}

		return nil
	})
}

func (f *ficsitCLI) RemoveMod(mod string) error {
	return f.action(ActionUninstall, newSimpleItem(mod), func(l *slog.Logger, taskUpdates chan<- taskUpdate) error {
		selectedInstallation := f.GetSelectedInstall()

		if selectedInstallation == nil {
			return fmt.Errorf("no installation selected")
		}

		l = l.With(
			slog.String("install", selectedInstallation.Path),
			slog.String("profile", selectedInstallation.Profile),
		)

		profile := f.GetProfile(selectedInstallation.Profile)

		profile.RemoveMod(mod)

		err := f.ficsitCli.Profiles.Save()
		if err != nil {
			l.Error("failed to save profile", slog.Any("error", err))
		}

		installErr := f.validateInstall(selectedInstallation, taskUpdates)

		if installErr != nil {
			l.Error("failed to install", slog.Any("error", installErr))
			return installErr
		}

		return nil
	})
}

func (f *ficsitCLI) EnableMod(mod string) error {
	return f.action(ActionEnable, newSimpleItem(mod), func(l *slog.Logger, taskUpdates chan<- taskUpdate) error {
		selectedInstallation := f.GetSelectedInstall()

		if selectedInstallation == nil {
			return fmt.Errorf("no installation selected")
		}

		l = l.With(
			slog.String("install", selectedInstallation.Path),
			slog.String("profile", selectedInstallation.Profile),
		)

		profile := f.GetProfile(selectedInstallation.Profile)

		profile.SetModEnabled(mod, true)

		err := f.ficsitCli.Profiles.Save()
		if err != nil {
			l.Error("failed to save profile", slog.Any("error", err))
		}

		installErr := f.validateInstall(selectedInstallation, taskUpdates)

		if installErr != nil {
			l.Error("failed to install", slog.Any("error", installErr))
			return installErr
		}

		return nil
	})
}

func (f *ficsitCLI) DisableMod(mod string) error {
	return f.action(ActionDisable, newSimpleItem(mod), func(l *slog.Logger, taskUpdates chan<- taskUpdate) error {
		selectedInstallation := f.GetSelectedInstall()

		if selectedInstallation == nil {
			return fmt.Errorf("no installation selected")
		}

		l = l.With(
			slog.String("install", selectedInstallation.Path),
			slog.String("profile", selectedInstallation.Profile),
		)

		profile := f.GetProfile(selectedInstallation.Profile)

		profile.SetModEnabled(mod, false)

		err := f.ficsitCli.Profiles.Save()
		if err != nil {
			l.Error("failed to save profile", slog.Any("error", err))
		}

		installErr := f.validateInstall(selectedInstallation, taskUpdates)

		if installErr != nil {
			l.Error("failed to install", slog.Any("error", installErr))
			return installErr
		}

		return nil
	})
}
