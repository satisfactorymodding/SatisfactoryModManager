package install_finders

type LegendaryGame struct {
	AppName           string   `json:"app_name"`
	BaseURLs          []string `json:"base_urls"`
	CanRunOffline     bool     `json:"can_run_offline"`
	EGL_GUID          string   `json:"egl_guid"`
	Executable        string   `json:"executable"`
	InstallPath       string   `json:"install_path"`
	InstallSize       int      `json:"install_size"`
	IsDLC             bool     `json:"is_dlc"`
	LaunchParameters  string   `json:"launch_parameters"`
	ManifestPath      string   `json:"manifest_path"`
	NeedsVerification bool     `json:"needs_verification"`
	RequiresOT        bool     `json:"requires_ot"`
	SavePath          string   `json:"save_path"`
	Title             string   `json:"title"`
	Version           string   `json:"version"`
}

type LegendaryData = map[string]LegendaryGame
