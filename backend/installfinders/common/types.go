package common

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

func (e InstallFindError) Cause() error {
	return e.Inner
}

type InstallFinderFunc func() ([]*Installation, []error)
