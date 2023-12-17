package installfinders

import (
	"path/filepath"
	"runtime"
	"strings"
)

type GameBranch string

var (
	BranchEarlyAccess  GameBranch = "Early Access"
	BranchExperimental GameBranch = "Experimental"
)

type InstallType string

var (
	InstallTypeWindowsClient InstallType = "WindowsClient"
	InstallTypeWindowsServer InstallType = "WindowsServer"
	InstallTypeLinuxServer   InstallType = "LinuxServer"
)

type LocationType string

var (
	LocationTypeLocal  LocationType = "Local"
	LocationTypeRemote LocationType = "Remote"
)

type Installation struct {
	Path       string       `json:"path"`
	Version    int          `json:"version"`
	Type       InstallType  `json:"type"`
	Location   LocationType `json:"location"`
	Branch     GameBranch   `json:"branch"`
	Launcher   string       `json:"launcher"`
	LaunchPath []string     `json:"launchPath"`
}

type InstallFindError struct {
	Inner error  `json:"cause"`
	Path  string `json:"path"`
}

func (e InstallFindError) Error() string {
	return e.Path + ": " + e.Inner.Error()
}

func (e InstallFindError) Causes() error {
	return e.Inner
}

type InstallFinderFunc func() ([]*Installation, []error)

func OsPathEqual(path1, path2 string) bool {
	path1 = filepath.Clean(path1)
	path2 = filepath.Clean(path2)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(path1, path2)
	}
	return path1 == path2
}

func FindAll(finders ...InstallFinderFunc) ([]*Installation, []error) {
	installs := make([]*Installation, 0)
	var errors []error
	for _, finder := range finders {
		foundInstalls, foundErrors := finder()
		for _, install := range foundInstalls {
			existing := false
			for i := range installs {
				if OsPathEqual(installs[i].Path, install.Path) {
					existing = true
					break
				}
			}
			if !existing {
				installs = append(installs, install)
			}
		}
		errors = append(errors, foundErrors...)
	}
	return installs, errors
}
