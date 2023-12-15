package apply

import "io"

type Apply interface {
	Apply(file io.Reader) error
	OnExit(restart bool) error
}
