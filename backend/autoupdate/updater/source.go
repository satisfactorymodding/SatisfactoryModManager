package updater

import "io"

type Source interface {
	GetLatestVersion() (string, error)
	GetChangelogs() (map[string]string, error)
	GetFile(filename string) (io.ReadCloser, int64, error)
}
