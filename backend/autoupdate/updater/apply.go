package updater

import "io"

type Apply interface {
	Apply(file io.Reader, checksum []byte) error
	OnExit(restart bool) error
}
