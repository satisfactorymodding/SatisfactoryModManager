package app

import (
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
)

type app struct {
	IsExpanded bool

	Restart bool

	stopSizeWatcher chan bool
}

var App = &app{}

func (a *app) CloseAndRestart() {
	a.Restart = true
	wailsRuntime.Quit(common.AppContext)
}
