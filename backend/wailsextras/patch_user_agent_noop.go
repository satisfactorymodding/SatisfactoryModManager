//go:build !windows && !cgo

// Cannot cross-platform lint the darwin and linux versions because CGO is disabled when cross-compiling

package wailsextras

func addUserAgent(_ string) {
	_ = allowUnexportedFieldAccess(getFrontendReflected()) // So that other functions are not marked as unused
	panic("this should never be reached")
}
