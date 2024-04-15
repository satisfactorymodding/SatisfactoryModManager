//go:build unix

package common

import (
	"path/filepath"
	"strings"
)

func WinePathProcessor(winePrefix string) func(string) string {
	return func(path string) string {
		return filepath.Join(winePrefix, "dosdevices", strings.ToLower(path[0:1])+strings.ReplaceAll(path[1:], "\\", "/"))
	}
}
