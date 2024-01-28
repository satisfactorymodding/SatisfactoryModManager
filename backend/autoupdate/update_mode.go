package autoupdate

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater"
)

type UpdateType struct {
	ArtifactName string
	Apply        updater.Apply
}

var updateMode = "none"

var updateTypes = map[string]func() UpdateType{}

func registerUpdateType(updateMode string, updateType func() UpdateType) {
	updateTypes[updateMode] = updateType
}

func getUpdateType() *UpdateType {
	getter := updateTypes[updateMode]
	if getter == nil {
		return nil
	}
	updateType := getter()
	return &updateType
}
