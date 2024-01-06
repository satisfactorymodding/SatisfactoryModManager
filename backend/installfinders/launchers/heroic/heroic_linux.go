package heroic

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func FindInstallations() ([]*common.Installation, []error) {
	return common.FindAll(findInstallationsNative, findInstallationsFlatpak, findInstallationsSnap)
}

func findInstallationsNative() ([]*common.Installation, []error) {
	return findInstallationsHeroic(false, "", "Heroic")
}

func findInstallationsFlatpak() ([]*common.Installation, []error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get user home dir")}
	}
	flatpakXdgConfigHome := filepath.Join(homeDir, ".var", "app", "com.heroicgameslauncher.hgl", "config")
	return findInstallationsHeroic(false, flatpakXdgConfigHome, "Heroic")
}

func findInstallationsSnap() ([]*common.Installation, []error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get user home dir")}
	}
	snapAppDir := filepath.Join(homeDir, "snap", "heroic")
	var latestSnapRevision int
	var latestSnapDirName string
	items, err := os.ReadDir(snapAppDir)
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to read heroic snap dir")}
	}
	for _, item := range items {
		if item.IsDir() {
			folderName := item.Name()
			var revision int
			if folderName[0] == 'x' {
				revision, err = strconv.Atoi(folderName[1:])
				if err != nil {
					continue
				}
			} else {
				revision, err = strconv.Atoi(folderName)
				if err != nil {
					continue
				}
			}
			if latestSnapDirName == "" || revision > latestSnapRevision {
				latestSnapRevision = revision
				latestSnapDirName = folderName
			}
		}
	}
	if latestSnapDirName == "" {
		return nil, []error{errors.New("no heroic snap folders found")}
	}

	return findInstallationsHeroic(true, filepath.Join(snapAppDir, latestSnapDirName, ".config"), "Heroic")
}
