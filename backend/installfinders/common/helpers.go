package common

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func OsPathEqual(path1, path2 string) bool {
	// Same file on disk (inode) means same install, regardless of path
	// spelling. This catches duplicates from case-only differences on
	// case-insensitive filesystems (e.g. macOS default APFS), where string
	// comparison alone sees two different installs.
	if info1, err1 := os.Stat(path1); err1 == nil {
		if info2, err2 := os.Stat(path2); err2 == nil && os.SameFile(info1, info2) {
			return true
		}
	}
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
