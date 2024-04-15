package heroic

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
)

func init() {
	launchers.Add("Heroic-snap", func() ([]*common.Installation, []error) {
		snapPath, err := getSnapPath()
		if err != nil {
			return nil, []error{fmt.Errorf("failed to get snap path: %w", err)}
		}

		return findInstallationsHeroic(true, filepath.Join(snapPath, ".config"), "Heroic")
	})
}

func getSnapPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}
	snapAppDir := filepath.Join(homeDir, "snap", "heroic")
	var latestSnapRevision int
	var latestSnapDirName string
	items, err := os.ReadDir(snapAppDir)
	if err != nil {
		return "", fmt.Errorf("failed to read heroic snap dir: %w", err)
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
		return "", fmt.Errorf("no heroic snap folders found")
	}
	return filepath.Join(snapAppDir, latestSnapDirName), nil
}
