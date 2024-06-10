package autoupdate

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/apply"
)

func init() {
	registerUpdateType("standalone", func() UpdateType {
		return UpdateType{
			ArtifactName: "SatisfactoryModManager_darwin_universal.zip",
			Apply: apply.MakeDarwinAppApply(apply.DarwinApplyConfig{
				AppName: "SatisfactoryModManager",
			}),
		}
	})
}
