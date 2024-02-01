package autoupdate

import (
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater"
)

type UpdateType struct {
	ArtifactName string
	Apply        updater.Apply
}

var updateTypes = map[string]func() UpdateType{}

func registerUpdateType(updateMode string, updateType func() UpdateType) {
	updateTypes[updateMode] = updateType
}

func shouldUseUpdater() bool {
	return viper.Get("update-mode") != "none"
}

func getUpdateType() *UpdateType {
	getter := updateTypes[viper.GetString("update-mode")]
	if getter == nil {
		return nil
	}
	updateType := getter()
	return &updateType
}
