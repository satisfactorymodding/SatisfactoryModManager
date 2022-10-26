package bindings

import (
	"context"

	"github.com/satisfactorymodding/SatisfactoryModManager/settings"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Settings struct {
	ctx context.Context
}

func MakeSettings() *Settings {
	s := &Settings{}
	return s
}

func (s *Settings) startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *Settings) FavouriteMod(modReference string) (bool, error) {
	idx := -1
	for i, mod := range settings.Settings.FavouriteMods {
		if mod == modReference {
			idx = i
			break
		}
	}
	if idx != -1 {
		return false, nil
	}
	settings.Settings.FavouriteMods = append(settings.Settings.FavouriteMods, modReference)
	err := settings.SaveSettings()
	if err != nil {
		return false, err
	}
	s.emitFavouriteMods()
	return true, nil
}

func (s *Settings) UnFavouriteMod(modReference string) bool {
	idx := -1
	for i, mod := range settings.Settings.FavouriteMods {
		if mod == modReference {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}
	settings.Settings.FavouriteMods = append(settings.Settings.FavouriteMods[:idx], settings.Settings.FavouriteMods[idx+1:]...)
	settings.SaveSettings()
	s.emitFavouriteMods()
	return true
}

func (s *Settings) GetFavouriteMods() []string {
	return settings.Settings.FavouriteMods
}

func (s *Settings) GetModFiltersOrder() string {
	return settings.Settings.ModFilters.Order
}

func (s *Settings) GetModFiltersFilter() string {
	return settings.Settings.ModFilters.Filter
}

func (s *Settings) SetModFiltersOrder(order string) {
	settings.Settings.ModFilters.Order = order
	settings.SaveSettings()
}

func (s *Settings) SetModFiltersFilter(filter string) {
	settings.Settings.ModFilters.Filter = filter
	settings.SaveSettings()
}

func (s *Settings) emitFavouriteMods() {
	wailsRuntime.EventsEmit(s.ctx, "favouriteMods", settings.Settings.FavouriteMods)
}

func (s *Settings) GetStartView() string {
	return settings.Settings.StartView
}

func (s *Settings) SetStartView(view string) {
	settings.Settings.StartView = view
	settings.SaveSettings()
}
