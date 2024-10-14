package common

import (
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
)

func OsPathEqual(path1, path2 string) bool {
	path1 = realPath(path1)
	path2 = realPath(path2)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(path1, path2)
	}
	return path1 == path2
}

func realPath(path string) string {
	newPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		slog.Warn("failed to evaluate symlink, using original path", slog.String("path", path), slog.Any("error", err))
		return path
	}
	return newPath
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
