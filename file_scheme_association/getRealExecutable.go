package file_scheme_association

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func getRealExecutablePath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", errors.Wrap(err, "failed to get executable path")
	}

	realExecutablePath, err := filepath.EvalSymlinks(executablePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve symlink")
	}

	return realExecutablePath, nil
}
