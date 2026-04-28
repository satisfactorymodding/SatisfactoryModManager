package ficsitcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/spf13/viper"
)

type ModpackReleaseResponse struct {
	GetModpackRelease struct {
		Lockfile string `json:"lockfile"`
	} `json:"getModpackRelease"`
}

type PlatformTarget struct {
	Hash string `json:"hash"`
	Link string `json:"link"`
}

type ModEntry struct {
	Dependencies interface{}               `json:"dependencies"`
	Targets      map[string]PlatformTarget `json:"targets"`
	Version      string                    `json:"version"`
}

type Lockfile struct {
	Mods    map[string]ModEntry `json:"mods"`
	Version int                 `json:"version"`
}

func (f *ficsitCLI) InstallModpackRelease(modpackID string, release string, name string) error {
	return f.action(ActionInstall, newItem(modpackID, release), func(l *slog.Logger, taskUpdates chan<- taskUpdate) error {
		selectedInstallation := f.GetSelectedInstall()
		if selectedInstallation == nil {
			return fmt.Errorf("no installation selected")
		}

		l = l.With(
			slog.String("install", selectedInstallation.Path),
			slog.String("profile", selectedInstallation.Profile),
		)
		f.AddProfile(name + "-" + release)
		profileErr := f.setProfileModpack(l, name+"-"+release)
		if profileErr != nil {
			l.Error("failed to set profile", slog.Any("error", profileErr))
			return fmt.Errorf("failed to set profile: %w", profileErr)
		}

		lockfile, err := getLockfile(modpackID, release)
		if err != nil {
			return fmt.Errorf("failed to get lockfile: %w", err)
		}

		for modID, mod := range lockfile.Mods {
			modErr := f.installModVersionModpack(l, modID, mod.Version)
			if modErr != nil {
				l.Error("failed to install mod",
					slog.String("mod", modID),
					slog.String("version", mod.Version),
					slog.Any("error", modErr))

				return fmt.Errorf("failed to install mod: %s@%s: %w",
					modID, mod.Version, modErr)
			}
		}

		installErr := f.apply(l, taskUpdates)
		if installErr != nil {
			l.Error("failed to install", slog.Any("error", installErr))
			return installErr
		}

		return nil
	})
}

type graphQLResponse struct {
	Data struct {
		GetModpackRelease struct {
			Lockfile string `json:"lockfile"`
		} `json:"getModpackRelease"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func getLockfile(modpackID string, release string) (Lockfile, error) {
	endpoint := viper.GetString("api-base") + viper.GetString("graphql-api")

	body := map[string]interface{}{
		"query": `query GetModpackRelease($modpackID: ModpackID!, $version: String!) {
			getModpackRelease(modpackID: $modpackID, version: $version) {
				lockfile
			}
		}`,
		"variables": map[string]interface{}{
			"modpackID": modpackID,
			"version":   release,
		},
	}

	jsonBody, _ := json.Marshal(body)
	resp, err := (&http.Client{}).Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return Lockfile{}, fmt.Errorf("failed to query GraphQL: %w", err)
	}
	defer resp.Body.Close()

	var gqlResp graphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return Lockfile{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return Lockfile{}, fmt.Errorf("GraphQL error: %s", gqlResp.Errors[0].Message)
	}
	if gqlResp.Data.GetModpackRelease.Lockfile == "" {
		return Lockfile{}, fmt.Errorf("lockfile not found for %s@%s", modpackID, release)
	}

	var lockfile Lockfile
	if err := json.Unmarshal([]byte(gqlResp.Data.GetModpackRelease.Lockfile), &lockfile); err != nil {
		return Lockfile{}, fmt.Errorf("failed to parse lockfile: %w", err)
	}
	return lockfile, nil
}

func (f *ficsitCLI) setProfileModpack(l *slog.Logger, profile string) error {
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
	f.EmitModsChange()

	return nil
}

func (f *ficsitCLI) installModVersionModpack(l *slog.Logger, mod string, version string) error {
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

	return nil
}
