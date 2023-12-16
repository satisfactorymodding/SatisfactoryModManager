package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/projectfile"
)

var (
	wailsJSONFilePath = "wails.json"
	viVersionFilePath = "build/windows/installer/vi_version.nsh"
)

func main() {
	wailsJSONFile, err := os.Open(wailsJSONFilePath)
	if err != nil {
		panic(err)
	}
	defer wailsJSONFile.Close()

	projectFile, err := io.ReadAll(wailsJSONFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(projectFile, &projectfile.ProjectFile)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(viVersionFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	version, _, _ := strings.Cut(projectfile.ProjectFile.Info.ProductVersion, "-")
	version, _, _ = strings.Cut(version, "+")

	for strings.Count(version, ".") < 3 {
		version += ".0"
	}

	_, _ = f.WriteString("# DO NOT EDIT - Generated automatically by build/windows/installer_version.go\n\n")
	_, _ = f.WriteString(fmt.Sprintf("!define VI_VERSION \"%s\"\n", version))
}
