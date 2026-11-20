package ficsitcli

import (
	"fmt"
	"log/slog"
	"maps"
	"sort"

	"github.com/satisfactorymodding/ficsit-cli/cli"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
)

// ProfileModDiff holds both profile snapshots, including identical mods.
// A nil side means the mod is absent from that profile.
type ProfileModDiff struct {
	Mod      string          `json:"mod"`
	Current  *cli.ProfileMod `json:"current"`
	Incoming *cli.ProfileMod `json:"incoming"`
}

// ProfileModSelection specifies the resulting enabled state and constraint in Current.
type ProfileModSelection struct {
	Mod      string          `json:"mod"`
	Current  *cli.ProfileMod `json:"current"`
	Incoming *cli.ProfileMod `json:"incoming"`
	Enabled  bool            `json:"enabled"`
	Version  string          `json:"version"`
}

type ProfileMergeResult struct {
	Applied    bool   `json:"applied"`
	ApplyError string `json:"applyError"`
}

func profileDiff(current, incoming *cli.Profile) []ProfileModDiff {
	references := make(map[string]bool, len(current.Mods)+len(incoming.Mods))
	for reference := range current.Mods {
		references[reference] = true
	}
	for reference := range incoming.Mods {
		references[reference] = true
	}

	changes := make([]ProfileModDiff, 0, len(references))
	for reference := range references {
		currentMod, inCurrent := current.Mods[reference]
		incomingMod, inIncoming := incoming.Mods[reference]
		change := ProfileModDiff{Mod: reference}
		if inCurrent {
			change.Current = &currentMod
		}
		if inIncoming {
			change.Incoming = &incomingMod
		}
		changes = append(changes, change)
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Mod < changes[j].Mod })
	return changes
}

func sameProfileMod(left, right *cli.ProfileMod) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

// mergeProfileMods validates the preview before building a separate mod map.
// Neither profile is changed, including when a selection is stale or invalid.
func mergeProfileMods(current, incoming *cli.Profile, selections []ProfileModSelection) (map[string]cli.ProfileMod, error) {
	if len(selections) == 0 {
		return nil, fmt.Errorf("no profile changes selected")
	}
	changes := profileDiff(current, incoming)
	byMod := make(map[string]ProfileModDiff, len(changes))
	for _, change := range changes {
		byMod[change.Mod] = change
	}
	seen := make(map[string]bool, len(selections))
	for _, selection := range selections {
		if seen[selection.Mod] {
			return nil, fmt.Errorf("duplicate change for mod %q", selection.Mod)
		}
		seen[selection.Mod] = true
		change, ok := byMod[selection.Mod]
		if !ok || !sameProfileMod(change.Current, selection.Current) || !sameProfileMod(change.Incoming, selection.Incoming) {
			return nil, fmt.Errorf("profiles changed for mod %q; refresh the comparison before merging", selection.Mod)
		}
		currentVersion := change.Current != nil && selection.Version == change.Current.Version
		incomingVersion := change.Incoming != nil && selection.Version == change.Incoming.Version
		if !currentVersion && !incomingVersion {
			return nil, fmt.Errorf("choose a current or incoming version for mod %q", selection.Mod)
		}
	}

	mods := maps.Clone(current.Mods)
	if mods == nil {
		mods = make(map[string]cli.ProfileMod)
	}
	for _, selection := range selections {
		if selection.Current == nil && !selection.Enabled {
			continue
		}
		mods[selection.Mod] = cli.ProfileMod{Version: selection.Version, Enabled: selection.Enabled}
	}
	return mods, nil
}

func (f *ficsitCLI) profilesForMerge(currentName, incomingName string) (*cli.Profile, *cli.Profile, error) {
	if currentName == incomingName {
		return nil, nil, fmt.Errorf("choose a different incoming profile")
	}
	current := f.GetProfile(currentName)
	if current == nil {
		return nil, nil, fmt.Errorf("current profile %q not found", currentName)
	}
	incoming := f.GetProfile(incomingName)
	if incoming == nil {
		return nil, nil, fmt.Errorf("incoming profile %q not found", incomingName)
	}
	return current, incoming, nil
}

func (f *ficsitCLI) GetProfileDiff(currentName, incomingName string) ([]ProfileModDiff, error) {
	if !f.actionMutex.TryLock() {
		return nil, fmt.Errorf("another operation in progress")
	}
	defer f.actionMutex.Unlock()

	current, incoming, err := f.profilesForMerge(currentName, incomingName)
	if err != nil {
		return nil, err
	}
	return profileDiff(current, incoming), nil
}

func (f *ficsitCLI) MergeProfile(currentName, incomingName string, selections []ProfileModSelection) (*ProfileMergeResult, error) {
	result := &ProfileMergeResult{}
	err := f.action(ActionMergeProfile, newSimpleItem(currentName), func(l *slog.Logger, taskChannel chan<- taskUpdate) error {
		if f.isGameRunning {
			return fmt.Errorf("cannot merge profiles while the game is running")
		}
		current, incoming, err := f.profilesForMerge(currentName, incomingName)
		if err != nil {
			return err
		}
		mods, err := mergeProfileMods(current, incoming, selections)
		if err != nil {
			return err
		}

		previousMods := current.Mods
		current.Mods = mods
		if err := f.ficsitCli.Profiles.Save(); err != nil {
			current.Mods = previousMods
			return fmt.Errorf("failed to save merged profile: %w", err)
		}
		f.EmitModsChange()
		if settings.Settings.QueueAutoStart {
			if err := f.applyProfile(currentName, l, taskChannel); err != nil {
				l.Error("merged profile saved, but failed to apply", slog.Any("error", err))
				result.ApplyError = err.Error()
			} else {
				result.Applied = true
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
