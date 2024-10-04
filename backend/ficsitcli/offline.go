package ficsitcli

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	ficsitcache "github.com/satisfactorymodding/ficsit-cli/cli/cache"
	"github.com/satisfactorymodding/ficsit-cli/cli/provider"
	resolver "github.com/satisfactorymodding/ficsit-resolver"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
)

func (f *ficsitCLI) GetOffline() bool {
	return f.ficsitCli.Provider.IsOffline()
}

func (f *ficsitCLI) SetOffline(offline bool) {
	f.ficsitCli.Provider.(*provider.MixedProvider).Offline = offline
	settings.Settings.Offline = offline
	_ = settings.SaveSettings()
}

type Mod struct {
	ModReference string                `json:"mod_reference"`
	Name         string                `json:"name"`
	Logo         *string               `json:"logo"` // Base64 encoded
	Authors      []string              `json:"authors"`
	Versions     []resolver.ModVersion `json:"versions"`
}

func (f *ficsitCLI) OfflineGetMods() ([]Mod, error) {
	cache, err := ficsitcache.GetCacheMods()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	mods := make([]Mod, 0)
	cache.Range(func(_ string, mod ficsitcache.Mod) bool {
		mods = append(mods, f.convertCacheFileToMod(mod))
		return true
	})
	return mods, nil
}

func (f *ficsitCLI) OfflineGetModsByReferences(modReferences []string) ([]Mod, error) {
	cache, err := ficsitcache.GetCacheMods()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	mods := make([]Mod, 0)
	cache.Range(func(modReference string, mod ficsitcache.Mod) bool {
		if !slices.Contains(modReferences, modReference) {
			return true
		}
		mods = append(mods, f.convertCacheFileToMod(mod))
		return true
	})
	return mods, nil
}

func (f *ficsitCLI) OfflineGetMod(modReference string) (Mod, error) {
	mod, err := ficsitcache.GetCacheMod(modReference)
	if err != nil {
		return Mod{}, fmt.Errorf("failed to get cache: %w", err)
	}
	return f.convertCacheFileToMod(mod), nil
}

func (f *ficsitCLI) convertCacheFileToMod(mod ficsitcache.Mod) Mod {
	if mod.ModReference == "" {
		return Mod{}
	}

	authors := make([]string, 0)

	for _, author := range strings.Split(mod.Author, ",") {
		authors = append(authors, strings.TrimSpace(author))
	}

	modVersions, err := f.ficsitCli.Provider.ModVersionsWithDependencies(context.TODO(), mod.ModReference)
	if err != nil {
		slog.Warn("failed to get mod versions", slog.String("mod", mod.ModReference), slog.Any("error", err))
	}

	return Mod{
		Name:         mod.Name,
		ModReference: mod.ModReference,
		Authors:      authors,
		Logo:         mod.Icon,
		Versions:     modVersions,
	}
}
