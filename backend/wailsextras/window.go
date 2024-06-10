package wailsextras

import (
	"context"
	"runtime"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

// WindowSetPosition wraps Wails's WindowSetPosition,
// but ensures that WindowSetPosition(WindowGetPosition())
// will result in the window remaining stationary
// regardless of OS's behaviour with multiple displays
func WindowSetPosition(ctx context.Context, x, y int) {
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		// WindowSetPosition expects relative to the current monitor,
		// but WindowGetPosition returns absolute
		curX, curY := wailsRuntime.WindowGetPosition(ctx)
		display := utils.GetDisplayBoundsAt(curX, curY)
		x -= display.Min.X
		y -= display.Min.Y
	}

	// It appears that on darwin Wails gets and sets the position
	// with values relative to the current monitor

	wailsRuntime.WindowSetPosition(ctx, x, y)
}
