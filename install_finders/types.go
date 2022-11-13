package install_finders

type GameBranch string

var (
	BRANCH_EARLY_ACCESS GameBranch = "Early Access"
	BRANCH_EXPERIMENTAL GameBranch = "Experimental"
)

type Installation struct {
	Path       string     `json:"path"`
	Version    int        `json:"version"`
	Branch     GameBranch `json:"branch"`
	Launcher   string     `json:"launcher"`
	LaunchPath []string   `json:"launchPath"`
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
