package install_finders

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
