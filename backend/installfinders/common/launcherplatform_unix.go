//go:build unix

package common

import (
	"path/filepath"
	"strings"
)

type winePlatform struct {
	winePrefix string
}

func WineLauncherPlatform(winePrefix string) Platform {
	return winePlatform{winePrefix: winePrefix}
}

func (p winePlatform) ProcessPath(path string) string {
	return filepath.Join(p.winePrefix, "dosdevices", strings.ToLower(path[0:1])+strings.ReplaceAll(path[1:], "\\", "/"))
}

func (p winePlatform) Os() string {
	return "windows"
}
