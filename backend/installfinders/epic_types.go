package installfinders

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
	CatalogItemID    string `json:"CatalogItemID"`
	ManifestLocation string `json:"ManifestLocation"`
	InstallationGUID string `json:"InstallationGUID"`
	MainGameAppName  string `json:"MainGameAppName"`
	AppVersionString string `json:"AppVersionString"`
	InstallLocation  string `json:"InstallLocation"`
}

type EpicGameManifest struct {
	AppName          string `json:"AppName"`
	CatalogNamespace string `json:"CatalogNamespace"`
	CatalogItemID    string `json:"CatalogItemID"`
}
