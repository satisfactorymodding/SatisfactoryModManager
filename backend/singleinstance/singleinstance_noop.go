//go:build bindings

package singleinstance

func RequestSingleInstanceLock() bool {
	return true
}

func ListenForSecondInstance() {
}
