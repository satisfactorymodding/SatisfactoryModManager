package install_finders

import (
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"
)

func FindInstallationsLinuxLutrisFlatpak() ([]*Installation, []error) {
	lutrisCmd := exec.Command("flatpak", "run", "net.lutris.Lutris", "-lj")
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
		currentInstalls, errs := findInstallationsWineEpic(lutrisGame.Directory, "Lutris - "+lutrisGame.Name, []string{"flatpak", "run", "net.lutris.Lutris", "lutris:rungame/" + lutrisGame.Slug})
		installs = append(installs, currentInstalls...)
		if errs != nil {
			findErrors = append(findErrors, errs...)
		}
	}
	return installs, findErrors
}
