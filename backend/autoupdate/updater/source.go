package updater

import "io"

type Source interface {
	GetLatestVersion(includePrereleases bool) (string, error)
	GetChangelogs() (map[string]string, error)
	GetFile(version string, filename string) (io.ReadCloser, int64, []byte, error)
}
