package autoupdate

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/apply"
)

func init() {
	registerUpdateType("standalone", func() UpdateType {
		return UpdateType{
			ArtifactName: "SatisfactoryModManager",
			Apply:        apply.MakeSingleFileApply(),
		}
	})
}
