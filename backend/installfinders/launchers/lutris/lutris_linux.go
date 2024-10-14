package lutris

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/epic"
)

type Game struct {
	ID        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Runner    string `json:"runner"`
	Directory string `json:"directory"`
}

func init() {
	launchers.Add("Lutris", func() ([]*common.Installation, []error) {
		return findInstallations([]string{"lutris"}, "Lutris")
	})
	launchers.Add("Lutris-flatpak", func() ([]*common.Installation, []error) {
		return findInstallations([]string{"flatpak", "run", "net.lutris.Lutris"}, "Lutris")
	})
}

func findInstallations(lutrisCmd []string, launcher string) ([]*common.Installation, []error) {
	lutrisLjCmd := makeLutrisCmd(lutrisCmd, "-lj")
	lutrisLj := exec.Command(lutrisLjCmd[0], lutrisLjCmd[1:]...)
	lutrisLj.Env = os.Environ()
	lutrisLj.Env = append(lutrisLj.Env, "LUTRIS_SKIP_INIT=1")
	if os.Getenv("APPIMAGE") != "" {
		// Must clear the APPIMAGE added entries of LD_LIBRARY_PATH, as well as PYTHONHOME and PYTHONPATH
		// otherwise lutris will load some appimage libraries and some system libraries, which are incompatible
		lutrisLj.Env = append(lutrisLj.Env, "PYTHONHOME=", "PYTHONPATH=")

		appdir := os.Getenv("APPDIR")
		ldLibraryPathEntries := strings.Split(os.Getenv("LD_LIBRARY_PATH"), ":")
		newLdLibraryPathEntries := []string{}
		for _, entry := range ldLibraryPathEntries {
			if !strings.HasPrefix(entry, appdir) {
				newLdLibraryPathEntries = append(newLdLibraryPathEntries, entry)
			}
		}
		lutrisLj.Env = append(lutrisLj.Env, "LD_LIBRARY_PATH="+strings.Join(newLdLibraryPathEntries, ":"))
	}
	outputBytes, err := lutrisLj.Output()
	if err != nil {
		return nil, []error{
			fmt.Errorf("failed to run lutris -lj: %w", err),
		}
	}
	var lutrisGames []Game
	err = json.Unmarshal(outputBytes, &lutrisGames)
	if err != nil {
		return nil, []error{
			fmt.Errorf("failed to parse lutris -lj output: %w", err),
		}
	}

	installs := []*common.Installation{}
	findErrors := []error{}
	for _, lutrisGame := range lutrisGames {
		currentInstalls, errs := epic.FindInstallationsWine(lutrisGame.Directory, launcher+" - "+lutrisGame.Name, makeLutrisCmd(lutrisCmd, "lutris:rungame/"+lutrisGame.Slug))
		installs = append(installs, currentInstalls...)
		if errs != nil {
			findErrors = append(findErrors, errs...)
		}
	}
	return installs, findErrors
}

func makeLutrisCmd(lutrisCmd []string, args ...string) []string {
	return append(lutrisCmd, args...)
}
