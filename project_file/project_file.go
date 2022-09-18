package project_file

type Info struct {
	CompanyName    string  `json:"companyName"`
	ProductName    string  `json:"productName"`
	ProductVersion string  `json:"productVersion"`
	Copyright      *string `json:"copyright"`
	Comments       *string `json:"comments"`
}

type Project struct {
	Name string `json:"name"`
	Info Info   `json:"info"`
}

var ProjectFile Project

func Version() string {
	return ProjectFile.Info.ProductVersion
}
