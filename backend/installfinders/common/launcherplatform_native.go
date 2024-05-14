package common

import (
	"runtime"
)

type nativePlatform struct{}

func NativePlatform() Platform {
	return nativePlatform{}
}

func (p nativePlatform) ProcessPath(path string) string {
	return path
}

func (p nativePlatform) Os() string {
	return runtime.GOOS
}
