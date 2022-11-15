package autoupdate

import (
	"path/filepath"
	"runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/autoupdate/updater/apply"
)

func getInstallType() (string, apply.Apply) {
	return "Satisfactory Mod Manager", apply.MakeSingleFileApply()
}
