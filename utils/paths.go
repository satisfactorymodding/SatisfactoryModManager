package utils

import (
	"os"

	"github.com/pkg/errors"
)

func EnsureDirExists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrapf(err, "failed to stat path %s", path)
		}

		err = os.MkdirAll(path, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to create directory %s", path)
		}
	}
	return nil
}
