package autoupdate

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/apply"
)

func init() {
	registerUpdateType("standalone", func() UpdateType {
		return UpdateType{
			ArtifactName: "SatisfactoryModManager_linux_amd64",
			Apply:        apply.MakeSingleFileApply(),
		}
	})
	registerUpdateType("appimage", func() UpdateType {
		return UpdateType{
			ArtifactName: "SatisfactoryModManager_linux_amd64.AppImage",
			Apply:        apply.MakeAppImageApply(),
		}
	})
}
