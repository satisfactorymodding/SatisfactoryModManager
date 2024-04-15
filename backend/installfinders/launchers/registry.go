package launchers

import "github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"

var finders map[string]common.InstallFinderFunc

func Add(id string, f common.InstallFinderFunc) {
	if finders == nil {
		finders = make(map[string]common.InstallFinderFunc)
	}
	if _, ok := finders[id]; ok {
		panic("launcher already registered")
	}
	finders[id] = f
}

func GetInstallFinders() map[string]common.InstallFinderFunc {
	return finders
}
