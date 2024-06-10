package common

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
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

func (p winePlatform) CacheDir() (string, error) {
	regCacheDir, err := p.getRegCacheDir()
	if err == nil {
		return regCacheDir, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return p.getDefaultCacheDir()
	}
	return "", err
}

func (p winePlatform) Os() string {
	return "windows"
}

func (p winePlatform) getRegCacheDir() (string, error) {
	userRegPath := filepath.Join(p.winePrefix, "user.reg")
	userRegBytes, err := os.ReadFile(userRegPath)
	if err != nil {
		return "", fmt.Errorf("failed to read user.reg: %w", err)
	}
	userRegText := string(userRegBytes)
	if strings.HasPrefix(userRegText, "WINE REGISTRY") {
		newLineIndex := strings.Index(userRegText, "\n")
		userRegText = userRegText[newLineIndex+1:]
	}
	userReg, err := ini.Load(strings.NewReader(userRegText))
	if err != nil {
		return "", fmt.Errorf("failed to load user.reg: %w", err)
	}
	return p.ProcessPath(userReg.Section(`Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Shell Folders`).Key("Local AppData").String()), nil
}

func (p winePlatform) getDefaultCacheDir() (string, error) {
	// Default can be either
	// modern: C:\Users\<linux username>\AppData\Local
	// legacy: C:\Users\<linux username>\Local Settings\Application Data
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	modernPath := p.ProcessPath(fmt.Sprintf("C:\\Users\\%s\\AppData\\Local", currentUser.Name))
	legacyPath := p.ProcessPath(fmt.Sprintf("C:\\Users\\%s\\Local Settings\\Application Data", currentUser.Name))

	if _, err := os.Stat(modernPath); err == nil {
		return modernPath, nil
	}

	return legacyPath, nil
}
