package updater

import "io"

type Source interface {
	GetLatestVersion(includePrereleases bool) (string, error)
	GetChangelogs() (map[string]string, error)
	GetFile(version string, filename string) (io.ReadCloser, int64, error)
}

type ChecksumSource interface {
	GetChecksumForFile(source Source, version string, filename string) ([]byte, error)
}

type Apply interface {
	Download(file io.Reader, checksum []byte) error
	Apply(restart bool) error
}
