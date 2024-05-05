package apply

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/minio/selfupdate"
)

type NsisApply struct {
	config NsisApplyConfig
}

type NsisApplyConfig struct {
	InstallerDownloadPath string
	IsAllUsers            bool
}

func MakeNsisApply(config NsisApplyConfig) *NsisApply {
	return &NsisApply{
		config: config,
	}
}

func (a *NsisApply) Download(file io.Reader, checksum []byte) error {
	// Use selfupdate to download and verify the update
	err := selfupdate.PrepareAndCheckBinary(file, selfupdate.Options{
		TargetPath: a.config.InstallerDownloadPath,
		Checksum:   checksum,
	})
	if err != nil {
		return fmt.Errorf("failed to download nsis update: %w", err)
	}

	// Variables as used by selfupdate
	updateDir := filepath.Dir(a.config.InstallerDownloadPath)
	filename := filepath.Base(a.config.InstallerDownloadPath)
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))

	// Ensure that the installer is actually at the expected path
	err = os.Rename(newPath, a.config.InstallerDownloadPath)
	if err != nil {
		return fmt.Errorf("failed to rename nsis update: %w", err)
	}
	return nil
}

func (a *NsisApply) Apply(restart bool) error {
	arguments := []string{"/S"}
	if a.config.IsAllUsers {
		arguments = append(arguments, "/AllUsers")
	} else {
		arguments = append(arguments, "/CurrentUser")
	}
	if restart {
		arguments = append(arguments, "/ForceRun")
	}
	cmd := exec.Command(a.config.InstallerDownloadPath, arguments...)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start nsis installer: %w", err)
	}
	err = cmd.Process.Release()
	if err != nil {
		return fmt.Errorf("failed to release nsis installer process: %w", err)
	}
	return nil
}
