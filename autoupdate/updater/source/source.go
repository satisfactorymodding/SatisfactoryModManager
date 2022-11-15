package source

import "io"

type Source interface {
	GetLatestVersion() (string, error)
	GetChangelogs() (map[string]string, error)
	GetFile() (io.ReadCloser, int64, error)
}
