package ficsitcli

import (
	"strings"

	"github.com/pkg/errors"
	ficsitcache "github.com/satisfactorymodding/ficsit-cli/cli/cache"
	"github.com/satisfactorymodding/ficsit-cli/cli/provider"
	"golang.org/x/exp/slices"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
)

func (f *FicsitCLI) GetOffline() bool {
	return f.ficsitCli.Provider.IsOffline()
}

func (f *FicsitCLI) SetOffline(offline bool) {
	f.ficsitCli.Provider.(*provider.MixedProvider).Offline = offline
	settings.Settings.Offline = offline
	_ = settings.SaveSettings()
}

type Mod struct {
	ModReference string       `json:"mod_reference"`
	Name         string       `json:"name"`
	Logo         *string      `json:"logo"` // Base64 encoded
	Authors      []string     `json:"authors"`
	Versions     []ModVersion `json:"versions"`
}

type ModVersion struct {
	Version      string                 `json:"version"`
	Size         int64                  `json:"size"`
	Dependencies []ModVersionDependency `json:"dependencies"`
}

type ModVersionDependency struct {
	ModReference string `json:"mod_id"`
	Condition    string `json:"condition"`
	Optional     bool   `json:"optional"`
}

func (f *FicsitCLI) OfflineGetMods() ([]Mod, error) {
	cache, err := ficsitcache.GetCache()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cache")
	}

	mods := make([]Mod, 0)
	cache.Range(func(modReference string, modFiles []ficsitcache.File) bool {
		if modReference == "SML" {
			return true
		}
		mods = append(mods, convertCacheFileToMod(modFiles))
		return true
	})
	return mods, nil
}

func (f *FicsitCLI) OfflineGetModsByReferences(modReferences []string) ([]Mod, error) {
	cache, err := ficsitcache.GetCache()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cache")
	}

	mods := make([]Mod, 0)
	cache.Range(func(modReference string, modFiles []ficsitcache.File) bool {
		if modReference == "SML" {
			return true
		}
		if !slices.Contains(modReferences, modReference) {
			return true
		}
		mods = append(mods, convertCacheFileToMod(modFiles))
		return true
	})
	return mods, nil
}

func (f *FicsitCLI) OfflineGetMod(modReference string) (Mod, error) {
	modFiles, err := ficsitcache.GetCacheMod(modReference)
	if err != nil {
		return Mod{}, errors.Wrap(err, "failed to get cache")
	}
	if modFiles == nil {
		return Mod{}, errors.New("mod not found")
	}
	return convertCacheFileToMod(modFiles), nil
}

type SMLVersion struct {
	Version             string `json:"version"`
	SatisfactoryVersion int    `json:"satisfactory_version"` // TODO
}

func (f *FicsitCLI) OfflineGetSMLVersions() ([]SMLVersion, error) {
	smlFiles, err := ficsitcache.GetCacheMod("SML")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cache")
	}
	if smlFiles == nil {
		return nil, errors.New("SML not found")
	}

	versions := make([]SMLVersion, 0)
	for _, file := range smlFiles {
		versions = append(versions, SMLVersion{
			Version:             file.Plugin.SemVersion,
			SatisfactoryVersion: 0,
		})
	}
	return versions, nil
}

func convertCacheFileToMod(files []ficsitcache.File) Mod {
	authors := make([]string, 0)

	for _, author := range strings.Split(files[0].Plugin.CreatedBy, ",") {
		authors = append(authors, strings.TrimSpace(author))
	}

	versions := make([]ModVersion, 0)
	for _, file := range files {
		dependencies := make([]ModVersionDependency, 0)

		for _, dependency := range file.Plugin.Plugins {
			if dependency.BasePlugin {
				continue
			}
			dependencies = append(dependencies, ModVersionDependency{
				ModReference: dependency.Name,
				Condition:    dependency.SemVersion,
				Optional:     dependency.Optional,
			})
		}

		versions = append(versions, ModVersion{
			Version:      file.Plugin.SemVersion,
			Size:         file.Size,
			Dependencies: dependencies,
		})
	}

	return Mod{
		Name:         files[0].Plugin.FriendlyName,
		ModReference: files[0].ModReference,
		Authors:      authors,
		Logo:         files[0].Icon,
		Versions:     versions,
	}
}
