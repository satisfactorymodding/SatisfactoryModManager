package app

import (
	"log/slog"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

func (a *app) WatchWindow() {
	a.stopSizeWatcher = make(chan bool)
	windowStateTicker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-a.stopSizeWatcher:
				return
			case <-windowStateTicker.C:
				if wailsRuntime.WindowIsMinimised(common.AppContext) {
					// When the window is minimized, the window position values are garbage
					continue
				}
				w, h := wailsRuntime.WindowGetSize(common.AppContext)
				x, y := wailsRuntime.WindowGetPosition(common.AppContext)
				changed := false
				if a.IsExpanded {
					if w != settings.Settings.ExpandedSize.Width {
						settings.Settings.ExpandedSize.Width = w
						changed = true
					}
				} else {
					if w != settings.Settings.UnexpandedSize.Width {
						settings.Settings.UnexpandedSize.Width = w
						changed = true
					}
				}
				if h != settings.Settings.ExpandedSize.Height {
					settings.Settings.ExpandedSize.Height = h
					changed = true
				}
				if h != settings.Settings.UnexpandedSize.Height {
					settings.Settings.UnexpandedSize.Height = h
					changed = true
				}
				if settings.Settings.WindowPosition == nil {
					settings.Settings.WindowPosition = &utils.Position{
						X: x,
						Y: y,
					}
					changed = true
				} else {
					if x != settings.Settings.WindowPosition.X {
						settings.Settings.WindowPosition.X = x
						changed = true
					}
					if y != settings.Settings.WindowPosition.Y {
						settings.Settings.WindowPosition.Y = y
						changed = true
					}
				}
				if changed {
					err := settings.SaveSettings()
					if err != nil {
						slog.Error("failed to save settings", slog.Any("error", err))
					}
				}
			}
		}
	}()
}

func (a *app) StopWindowWatcher() {
	if a.stopSizeWatcher != nil {
		close(a.stopSizeWatcher)
	}
}
