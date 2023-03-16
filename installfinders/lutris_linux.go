package installfinders

import (
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"
)

func FindInstallationsLinuxLutris() ([]*Installation, []error) {
	lutrisCmd := exec.Command("lutris", "-lj")
	outputBytes, err := lutrisCmd.Output()
	if err != nil {
		return nil, []error{
			errors.Wrap(err, "failed to run lutris -lj"),
		}
	}
	var lutrisGames []LutrisGame
	err = json.Unmarshal(outputBytes, &lutrisGames)
	if err != nil {
		return nil, []error{
			errors.Wrap(err, "failed to parse lutris -lj output"),
		}
	}

	installs := []*Installation{}
	findErrors := []error{}
	for _, lutrisGame := range lutrisGames {
		currentInstalls, errs := findInstallationsWineEpic(lutrisGame.Directory, "Lutris - "+lutrisGame.Name, []string{"lutris", "lutris:rungame/" + lutrisGame.Slug})
		installs = append(installs, currentInstalls...)
		if errs != nil {
			findErrors = append(findErrors, errs...)
		}
	}
	return installs, findErrors
}
