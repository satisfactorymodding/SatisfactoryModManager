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

func (s *Settings) FavoriteMod(modReference string) (bool, error) {
	idx := -1
	for i, mod := range settings.Settings.FavoriteMods {
		if mod == modReference {
			idx = i
			break
		}
	}
	if idx != -1 {
		return false, nil
	}
	settings.Settings.FavoriteMods = append(settings.Settings.FavoriteMods, modReference)
	err := settings.SaveSettings()
	if err != nil {
		return false, err
	}
	s.emitFavoriteMods()
	return true, nil
}

func (s *Settings) UnFavoriteMod(modReference string) bool {
	idx := -1
	for i, mod := range settings.Settings.FavoriteMods {
		if mod == modReference {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}
	settings.Settings.FavoriteMods = append(settings.Settings.FavoriteMods[:idx], settings.Settings.FavoriteMods[idx+1:]...)
	settings.SaveSettings()
	s.emitFavoriteMods()
	return true
}

func (s *Settings) GetFavoriteMods() []string {
	return settings.Settings.FavoriteMods
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

func (s *Settings) emitFavoriteMods() {
	wailsRuntime.EventsEmit(s.ctx, "favoriteMods", settings.Settings.FavoriteMods)
}

func (s *Settings) GetStartView() settings.View {
	return settings.Settings.StartView
}

func (s *Settings) SetStartView(view settings.View) {
	settings.Settings.StartView = view
	settings.SaveSettings()
}

func (s *Settings) GetKonami() bool {
	return settings.Settings.Konami
}

func (s *Settings) SetKonami(value bool) {
	settings.Settings.Konami = value
	settings.SaveSettings()
}

func (s *Settings) GetLaunchButton() string {
	return settings.Settings.LaunchButton
}

func (s *Settings) SetLaunchButton(value string) {
	settings.Settings.LaunchButton = value
	settings.SaveSettings()
}

func (s *Settings) GetQueueAutoStart() bool {
	return settings.Settings.QueueAutoStart
}

func (s *Settings) SetQueueAutoStart(value bool) {
	settings.Settings.QueueAutoStart = value
	settings.SaveSettings()
}
