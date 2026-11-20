package ficsitcli

import (
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/puzpuzpuz/xsync/v3"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	resolver "github.com/satisfactorymodding/ficsit-resolver"
)

func TestProfileDiff(t *testing.T) {
	current := &cli.Profile{Mods: map[string]cli.ProfileMod{
		"Same":    {Version: ">=0.0.0", Enabled: true},
		"Version": {Version: "1.0.0", Enabled: true},
		"State":   {Version: "2.0.0", Enabled: true},
		"Removed": {Version: "3.0.0", Enabled: false},
	}}
	incoming := &cli.Profile{Mods: map[string]cli.ProfileMod{
		"Same":    {Version: ">=0.0.0", Enabled: true},
		"Version": {Version: "1.1.0", Enabled: true},
		"State":   {Version: "2.0.0", Enabled: false},
		"Added":   {Version: "4.0.0", Enabled: false},
	}}
	changes := profileDiff(current, incoming)
	if len(changes) != 5 {
		t.Fatalf("expected all 5 mods, got %#v", changes)
	}
	if changes[0].Mod != "Added" || changes[1].Mod != "Removed" || changes[2].Mod != "Same" || changes[3].Mod != "State" || changes[4].Mod != "Version" {
		t.Fatalf("differences are not sorted: %#v", changes)
	}
	if changes[0].Current != nil || *changes[0].Incoming != incoming.Mods["Added"] {
		t.Errorf("incorrect addition: %#v", changes[0])
	}
	if changes[1].Incoming != nil || *changes[1].Current != current.Mods["Removed"] {
		t.Errorf("incorrect removal: %#v", changes[1])
	}
	if !changes[2].Current.Enabled || !changes[2].Incoming.Enabled {
		t.Errorf("mutually enabled mod was lost: %#v", changes[2])
	}
	if !changes[3].Current.Enabled || changes[3].Incoming.Enabled {
		t.Errorf("enabled-state change was lost: %#v", changes[3])
	}
	if changes[4].Current.Version != "1.0.0" || changes[4].Incoming.Version != "1.1.0" {
		t.Errorf("version change was lost: %#v", changes[4])
	}
	// A preview must be a snapshot rather than a pointer into profile state.
	changes[0].Incoming.Version = "changed"
	if incoming.Mods["Added"].Version != "4.0.0" {
		t.Fatal("preview changed Incoming")
	}
}

func TestProfileDiffEmptyAndIdentical(t *testing.T) {
	for _, profile := range []*cli.Profile{
		{},
		{Mods: map[string]cli.ProfileMod{}},
	} {
		if changes := profileDiff(profile, profile); changes == nil || len(changes) != 0 {
			t.Fatalf("expected an empty JSON array, got %#v", changes)
		}
	}
	identical := &cli.Profile{Mods: map[string]cli.ProfileMod{"Mod": {Version: "1.0.0", Enabled: false}}}
	if entries := profileDiff(identical, identical); len(entries) != 1 || entries[0].Current == nil || entries[0].Incoming == nil {
		t.Fatalf("identical disabled mods must remain visible: %#v", entries)
	}
	// Presence differs even if the stored mod has zero-valued fields.
	changes := profileDiff(&cli.Profile{}, &cli.Profile{Mods: map[string]cli.ProfileMod{"Zero": {}}})
	if len(changes) != 1 || changes[0].Current != nil || changes[0].Incoming == nil {
		t.Fatalf("missing presence difference: %#v", changes)
	}
}

func TestMergeProfileModsSelective(t *testing.T) {
	current := &cli.Profile{
		Name:            "Current",
		RequiredTargets: []resolver.TargetName{"Windows"},
		Mods: map[string]cli.ProfileMod{
			"Keep":     {Version: "1.0.0", Enabled: true},
			"Change":   {Version: "2.0.0", Enabled: true},
			"Unpicked": {Version: "3.0.0", Enabled: true},
		},
	}
	incoming := &cli.Profile{Mods: map[string]cli.ProfileMod{
		"Added":    {Version: ">=4.0.0", Enabled: false},
		"Change":   {Version: "2.1.0", Enabled: false},
		"Unpicked": {Version: "3.1.0", Enabled: false},
	}}
	beforeCurrent, beforeIncoming := maps.Clone(current.Mods), maps.Clone(incoming.Mods)
	var selected []ProfileModSelection
	for _, change := range profileDiff(current, incoming) {
		if change.Mod == "Added" || change.Mod == "Change" {
			selected = append(selected, profileSelection(change, change.Mod == "Added", change.Incoming.Version))
		}
	}
	merged, err := mergeProfileMods(current, incoming, selected)
	if err != nil {
		t.Fatal(err)
	}
	want := maps.Clone(beforeCurrent)
	want["Added"] = cli.ProfileMod{Version: incoming.Mods["Added"].Version, Enabled: true}
	want["Change"] = incoming.Mods["Change"]
	if !maps.Equal(merged, want) {
		t.Fatalf("unexpected result: got %#v, want %#v", merged, want)
	}
	merged["Change"] = cli.ProfileMod{}
	if !maps.Equal(current.Mods, beforeCurrent) || !maps.Equal(incoming.Mods, beforeIncoming) {
		t.Fatal("merge mutated an input profile")
	}
	if current.Name != "Current" || !reflect.DeepEqual(current.RequiredTargets, []resolver.TargetName{"Windows"}) {
		t.Fatal("merge changed profile metadata")
	}
}

func TestMergeProfileModsDisablesCurrentOnlyMod(t *testing.T) {
	current := &cli.Profile{Mods: map[string]cli.ProfileMod{"Remove": {Version: "1.0.0", Enabled: true}}}
	incoming := &cli.Profile{}
	entry := profileDiff(current, incoming)[0]
	merged, err := mergeProfileMods(current, incoming, []ProfileModSelection{profileSelection(entry, false, "1.0.0")})
	if err != nil || len(merged) != 1 || merged["Remove"].Enabled || merged["Remove"].Version != "1.0.0" {
		t.Fatalf("unchecking must disable the mod and preserve its constraint: %#v, %v", merged, err)
	}
	if !current.Mods["Remove"].Enabled {
		t.Fatal("disabling mutated Current before saving")
	}
}

func TestMergeProfileModsIntoEmptyProfile(t *testing.T) {
	current := &cli.Profile{}
	incoming := &cli.Profile{Mods: map[string]cli.ProfileMod{"Added": {Version: "1.0.0", Enabled: false}}}
	entry := profileDiff(current, incoming)[0]
	merged, err := mergeProfileMods(current, incoming, []ProfileModSelection{profileSelection(entry, true, "1.0.0")})
	if err != nil || len(merged) != 1 || !merged["Added"].Enabled || merged["Added"].Version != "1.0.0" {
		t.Fatalf("merge into empty profile failed: %#v, %v", merged, err)
	}
	if incoming.Mods["Added"].Enabled {
		t.Fatal("checking an Incoming-only mod changed Incoming")
	}
}

func profileSelection(entry ProfileModDiff, enabled bool, version string) ProfileModSelection {
	return ProfileModSelection{Mod: entry.Mod, Current: entry.Current, Incoming: entry.Incoming, Enabled: enabled, Version: version}
}

func TestMergeProfileModsDisablesAllIncludingMutuallyEnabled(t *testing.T) {
	current := &cli.Profile{Mods: map[string]cli.ProfileMod{
		"Common": {Version: "1.0.0", Enabled: true},
		"Local":  {Version: "2.0.0", Enabled: true},
	}}
	incoming := &cli.Profile{Mods: map[string]cli.ProfileMod{"Common": {Version: "1.0.0", Enabled: true}}}
	var selected []ProfileModSelection
	for _, entry := range profileDiff(current, incoming) {
		selected = append(selected, profileSelection(entry, false, entry.Current.Version))
	}
	merged, err := mergeProfileMods(current, incoming, selected)
	if err != nil || len(merged) != 2 || merged["Common"].Enabled || merged["Local"].Enabled {
		t.Fatalf("must allow saving with no enabled mods: %#v, %v", merged, err)
	}
}

func TestMergeProfileModsVersionChoiceIsIndependentOfEnabled(t *testing.T) {
	current := &cli.Profile{Mods: map[string]cli.ProfileMod{"Mod": {Version: "1.0.0", Enabled: true}}}
	incoming := &cli.Profile{Mods: map[string]cli.ProfileMod{"Mod": {Version: "2.0.0", Enabled: false}}}
	entry := profileDiff(current, incoming)[0]
	for _, version := range []string{"1.0.0", "2.0.0"} {
		merged, err := mergeProfileMods(current, incoming, []ProfileModSelection{profileSelection(entry, true, version)})
		if err != nil || !merged["Mod"].Enabled || merged["Mod"].Version != version {
			t.Fatalf("enabled state must not be copied from Incoming: %#v, %v", merged, err)
		}
	}
}

func TestMergeProfileModsUncheckedIncomingOnlyStaysAbsent(t *testing.T) {
	current := &cli.Profile{}
	incoming := &cli.Profile{Mods: map[string]cli.ProfileMod{"Mod": {Version: "1.0.0", Enabled: true}}}
	entry := profileDiff(current, incoming)[0]
	merged, err := mergeProfileMods(current, incoming, []ProfileModSelection{profileSelection(entry, false, "1.0.0")})
	if err != nil || len(merged) != 0 {
		t.Fatalf("an unchecked Incoming-only mod must stay absent: %#v, %v", merged, err)
	}
}

func TestMergeProfileModsRejectsInvalidSelections(t *testing.T) {
	for _, scenario := range []string{"empty", "duplicate", "unknown", "current changed", "incoming changed", "incoming removed", "forged version", "forged snapshot", "forged absence"} {
		t.Run(scenario, func(t *testing.T) {
			current := &cli.Profile{Mods: map[string]cli.ProfileMod{"Mod": {Version: "1.0.0", Enabled: true}}}
			incoming := &cli.Profile{Mods: map[string]cli.ProfileMod{"Mod": {Version: "2.0.0", Enabled: false}}}
			selected := []ProfileModSelection{profileSelection(profileDiff(current, incoming)[0], true, "2.0.0")}
			switch scenario {
			case "empty":
				selected = nil
			case "duplicate":
				selected = append(selected, selected[0])
			case "unknown":
				selected[0].Mod = "Unknown"
			case "current changed":
				current.Mods["Mod"] = cli.ProfileMod{Version: "3.0.0", Enabled: true}
			case "incoming changed":
				incoming.Mods["Mod"] = cli.ProfileMod{Version: "3.0.0", Enabled: false}
			case "incoming removed":
				delete(incoming.Mods, "Mod")
			case "forged version":
				selected[0].Version = "9.0.0"
			case "forged snapshot":
				selected[0].Incoming.Version = "9.0.0"
			case "forged absence":
				selected[0].Incoming = nil
			}
			beforeCurrent, beforeIncoming := maps.Clone(current.Mods), maps.Clone(incoming.Mods)
			if merged, err := mergeProfileMods(current, incoming, selected); err == nil || merged != nil {
				t.Fatalf("accepted invalid selection: %#v, %v", merged, err)
			}
			if !maps.Equal(current.Mods, beforeCurrent) || !maps.Equal(incoming.Mods, beforeIncoming) {
				t.Fatal("rejected merge mutated a profile")
			}
		})
	}
}

func TestProfilesForMergeValidation(t *testing.T) {
	f := &ficsitCLI{ficsitCli: &cli.GlobalContext{Profiles: &cli.Profiles{
		Profiles: map[string]*cli.Profile{"Current": {Name: "Current"}, "Incoming": {Name: "Incoming"}},
	}}}
	for _, pair := range [][2]string{{"Current", "Current"}, {"Missing", "Incoming"}, {"Current", "Missing"}} {
		if _, _, err := f.profilesForMerge(pair[0], pair[1]); err == nil {
			t.Fatalf("accepted invalid profile pair: %v", pair)
		}
	}
	if _, _, err := f.profilesForMerge("Current", "Incoming"); err != nil {
		t.Fatal(err)
	}
}

func profileMergeTestInstallation(t *testing.T, profile string, vanilla bool) *cli.Installation {
	t.Helper()
	root := t.TempDir()
	versionDir := filepath.Join(root, "Engine", "Binaries", "Win64")
	if err := os.MkdirAll(versionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{filepath.Join(root, "FactoryGame.exe"), filepath.Join(versionDir, "FactoryGame-Win64-Shipping.version")} {
		if err := os.WriteFile(file, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return &cli.Installation{Path: root, Profile: profile, Vanilla: vanilla}
}

func TestChosenMergeProfileTargetsItsOwnInstallations(t *testing.T) {
	selected := profileMergeTestInstallation(t, "Selected", true)
	destination := profileMergeTestInstallation(t, "Destination", false)
	// These paths deliberately do not exist: excluded installs must not be inspected.
	inactive := &cli.Installation{Path: "inactive", Profile: "Destination", Vanilla: true}
	unrelated := &cli.Installation{Path: "unrelated", Profile: "Other"}
	f := &ficsitCLI{
		ficsitCli: &cli.GlobalContext{
			Profiles: &cli.Profiles{Profiles: map[string]*cli.Profile{
				"Selected":    {Name: "Selected"},
				"Destination": {Name: "Destination"},
				"Other":       {Name: "Other"},
				"Unused":      {Name: "Unused", RequiredTargets: []resolver.TargetName{"Windows"}},
			}},
			Installations: &cli.Installations{
				SelectedInstallation: selected.Path,
				Installations:        []*cli.Installation{selected, destination, inactive, unrelated},
			},
		},
		installationMetadata: xsync.NewMapOf[string, installationMetadata](),
	}
	for _, install := range f.ficsitCli.Installations.Installations {
		f.installationMetadata.Store(install.Path, installationMetadata{State: InstallStateValid})
	}
	installs, profile, err := f.getInstallsToApplyForProfile("Destination")
	if err != nil {
		t.Fatal(err)
	}
	if profile != f.GetProfile("Destination") || len(installs) != 1 || installs[0].install != destination || installs[0].targetName != "Windows" {
		t.Fatalf("wrong merge destination installations: %#v, %#v", profile, installs)
	}
	if f.GetSelectedInstall() != selected || selected.Profile != "Selected" {
		t.Fatal("choosing a merge destination changed the selected installation")
	}
	// The ordinary Apply path still includes the selected install when mods are off.
	installs, profile, err = f.getInstallsToApply()
	if err != nil || profile != f.GetProfile("Selected") || len(installs) != 1 || installs[0].install != selected {
		t.Fatalf("selected-profile apply changed: %#v, %v", installs, err)
	}
	if _, _, err := f.getInstallsToApplyForProfile("Missing"); err == nil {
		t.Fatal("accepted a missing destination")
	}
	if err := f.applyProfile("Unused", slog.Default(), nil); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(f.GetProfile("Unused").RequiredTargets, []resolver.TargetName{"Windows"}) {
		t.Fatal("applying an unused destination changed its target constraints")
	}
}
