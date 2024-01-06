package autoupdate

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater/apply"
)

func getInstallType() (string, apply.Apply) {
	return "Satisfactory Mod Manager", apply.MakeSingleFileApply()
}
