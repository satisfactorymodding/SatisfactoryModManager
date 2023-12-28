package common

import (
	"path/filepath"
	"runtime"
	"strings"
)

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
