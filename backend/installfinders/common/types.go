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
	SavedPath  string       `json:"-"`
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

var AllInstallTypes = []struct {
	Value  InstallType
	TSName string
}{
	{InstallTypeWindowsClient, "WINDOWS"},
	{InstallTypeWindowsServer, "WINDOWS_SERVER"},
	{InstallTypeLinuxServer, "LINUX_SERVER"},
}

var AllBranches = []struct {
	Value  GameBranch
	TSName string
}{
	{BranchEarlyAccess, "EARLY_ACCESS"},
	{BranchExperimental, "EXPERIMENTAL"},
}

var AllLocationTypes = []struct {
	Value  LocationType
	TSName string
}{
	{LocationTypeLocal, "LOCAL"},
	{LocationTypeRemote, "REMOTE"},
}
