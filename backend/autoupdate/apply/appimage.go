package apply

import (
	"fmt"
	"io"
	"os"

	"github.com/minio/selfupdate"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

type AppImageApply struct{}

func MakeAppImageApply() *AppImageApply {
	return &AppImageApply{}
}

func (a *AppImageApply) Download(file io.Reader, checksum []byte) error {
	appimagePath := os.Getenv("APPIMAGE")
	if appimagePath == "" {
		return fmt.Errorf("APPIMAGE environment variable not set, executable not an appimage")
	}

	err := selfupdate.PrepareAndCheckBinary(file, selfupdate.Options{
		Checksum:   checksum,
		TargetPath: appimagePath,
	})
	if err != nil {
		return fmt.Errorf("failed to download AppImage update: %w", err)
	}
	return nil
}

func (a *AppImageApply) Apply(restart bool) error {
	appimagePath := os.Getenv("APPIMAGE")
	if appimagePath == "" {
		return fmt.Errorf("APPIMAGE environment variable not set, executable not an appimage")
	}

	err := selfupdate.CommitBinary(selfupdate.Options{
		TargetPath: appimagePath,
	})
	if err != nil {
		return fmt.Errorf("failed to commit AppImage update: %w", err)
	}
	if restart {
		err := utils.Restart()
		if err != nil {
			return fmt.Errorf("failed to relaunch after update: %w", err)
		}
	}
	return nil
}
