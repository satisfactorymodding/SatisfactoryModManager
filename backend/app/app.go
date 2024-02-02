package app

type app struct {
	IsExpanded bool

	stopSizeWatcher chan bool
}

var App = &app{}
