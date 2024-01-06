package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func IsIn(dir, path string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}
	return filepath.IsLocal(rel)
}

func CopyRecursive(from, to string) error {
	return filepath.Walk(from, func(path string, info os.FileInfo, err error) error { //nolint:wrapcheck
		if err != nil {
			return err
		}
		if IsIn(to, path) {
			return nil
		}
		relPath, err := filepath.Rel(from, path)
		if err != nil {
			return err //nolint:wrapcheck
		}
		newPath := filepath.Join(to, relPath)
		if info.IsDir() {
			err := os.Mkdir(newPath, 0o755)
			if err != nil && !os.IsExist(err) {
				return err //nolint:wrapcheck
			}
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err //nolint:wrapcheck
		}
		defer f.Close()
		f2, err := os.Create(newPath)
		if err != nil {
			return err //nolint:wrapcheck
		}
		defer f2.Close()
		_, err = io.Copy(f2, f)
		return err //nolint:wrapcheck
	})
}

func MoveRecursive(from, to string) (bool, error) {
	err := CopyRecursive(from, to)
	if err != nil {
		return false, errors.Wrapf(err, "failed to copy %s to %s", from, to)
	}
	err = filepath.Walk(from, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			return nil
		}
		if IsIn(path, to) {
			// Skip parent directories of destination
			return nil
		}
		if IsIn(to, path) {
			// Skip contents of destination
			return nil
		}
		err = os.RemoveAll(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return err //nolint:wrapcheck
			}
		}
		return nil
	})
	if err != nil {
		return true, errors.Wrapf(err, "failed to remove %s", from)
	}
	return true, nil
}
