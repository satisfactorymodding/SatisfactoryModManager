package install_finders

type Installation struct {
	Path       string   `json:"path"`
	Version    int      `json:"version"`
	Branch     string   `json:"branch"`
	Launcher   string   `json:"launcher"`
	LaunchPath []string `json:"launchPath"`
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

type GameVersionFile struct {
	MajorVersion         int    `json:"MajorVersion"`
	MinorVersion         int    `json:"MinorVersion"`
	PatchVersion         int    `json:"PatchVersion"`
	Changelist           int    `json:"Changelist"`
	CompatibleChangelist int    `json:"CompatibleChangelist"`
	IsLicenseeVersion    int    `json:"IsLicenseeVersion"`
	IsPromotedBuild      int    `json:"IsPromotedBuild"`
	BranchName           string `json:"BranchName"`
	BuildID              string `json:"BuildId"`
}

type EpicManifest struct {
	CatalogNamespace string `json:"CatalogNamespace"`
	CatalogItemId    string `json:"CatalogItemId"`
	ManifestLocation string `json:"ManifestLocation"`
	InstallationGuid string `json:"InstallationGuid"`
	MainGameAppName  string `json:"MainGameAppName"`
	AppVersionString string `json:"AppVersionString"`
	InstallLocation  string `json:"InstallLocation"`
}

type EpicGameManifest struct {
	AppName          string `json:"AppName"`
	CatalogNamespace string `json:"CatalogNamespace"`
	CatalogItemId    string `json:"CatalogItemId"`
}
