package installfinders

import (
	"log/slog"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"

	_ "github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/all" // register all launchers
)

func FindInstallations() ([]*common.Installation, []error) {
	registrations := launchers.GetInstallFinders()

	slog.Debug("finding installations", slog.String("launchers", strings.Join(maps.Keys(registrations), ",")))

	return common.FindAll(maps.Values(registrations)...)
}
