package app

import (
	"image"
	"log/slog"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

func (a *app) WatchWindow() {
	a.stopSizeWatcher = make(chan bool)
	windowStateTicker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-a.stopSizeWatcher:
				return
			case <-windowStateTicker.C:
				a.ensureWindowVisible()
				a.saveWindowState()
			}
		}
	}()
}

func (a *app) ensureWindowVisible() {
	if wailsRuntime.WindowIsMinimised(common.AppContext) {
		// When the window is minimized, the window position values are garbage
		return
	}
	x, y := wailsRuntime.WindowGetPosition(common.AppContext)
	w, h := wailsRuntime.WindowGetSize(common.AppContext)

	window := image.Rect(x, y, x+w, y+h)

	displays := utils.GetDisplayBounds()

	if len(displays) == 0 {
		slog.Warn("no displays found, cannot check if window is reachable")
		return
	}

	for _, display := range displays {
		if display.Overlaps(window) {
			return
		}
		if display.Dx() == 0 || display.Dy() == 0 {
			slog.Warn("display has no size", slog.Any("display", display))
			return
		}
	}

	// If the window is not visible, move it to the center of the primary display
	wailsRuntime.WindowCenter(common.AppContext)
}

func (a *app) saveWindowState() {
	if wailsRuntime.WindowIsMinimised(common.AppContext) {
		// When the window is minimized, the window position values are garbage
		return
	}
	w, h := wailsRuntime.WindowGetSize(common.AppContext)
	x, y := wailsRuntime.WindowGetPosition(common.AppContext)
	maximized := wailsRuntime.WindowIsMaximised(common.AppContext)
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
	if maximized != settings.Settings.Maximized {
		settings.Settings.Maximized = maximized
		changed = true
	}
	if changed {
		err := settings.SaveSettings()
		if err != nil {
			slog.Error("failed to save settings", slog.Any("error", err))
		}
	}
}

func (a *app) StopWindowWatcher() {
	if a.stopSizeWatcher != nil {
		close(a.stopSizeWatcher)
	}
}
