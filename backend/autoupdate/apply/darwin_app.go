package apply

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/minio/selfupdate"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

type DarwinAppApply struct {
	config DarwinApplyConfig
}

type DarwinApplyConfig struct {
	AppName string
}

func MakeDarwinAppApply(config DarwinApplyConfig) *DarwinAppApply {
	return &DarwinAppApply{
		config: config,
	}
}

func (a *DarwinAppApply) Apply(file io.Reader, checksum []byte) error {
	executable, err := exec.LookPath(os.Args[0])
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// The archive contains a .app directory, so we need to find the .app directory to replace
	appDir := executable
	for !strings.HasSuffix(filepath.Base(appDir), ".app") && appDir != filepath.Dir(appDir) {
		appDir = filepath.Dir(appDir)
	}

	if !strings.HasSuffix(filepath.Base(appDir), ".app") {
		return fmt.Errorf("failed to find .app directory")
	}

	options := selfupdate.Options{
		Checksum:   checksum,
		TargetPath: appDir,
	}

	// Variables as used by selfupdate
	updateDir := filepath.Dir(options.TargetPath)
	filename := filepath.Base(options.TargetPath)
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))

	// Ensure that newPath does not exist, since it might be a dir
	err = os.RemoveAll(newPath)
	if err != nil {
		return fmt.Errorf("failed to remove existing darwin update: %w", err)
	}

	err = selfupdate.PrepareAndCheckBinary(file, options)
	if err != nil {
		return fmt.Errorf("failed to save darwin update: %w", err)
	}

	// Now we should have a newPath file that is a zip
	// which contains only the ${config.AppName}.app dir
	// We need to unzip it, but give the ${config.AppName}.app dir
	// the newPath name

	tmpDir := fmt.Sprintf("%s.tmp", newPath)

	// Remove the ${tmpDir} dir
	err = os.RemoveAll(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to remove darwin update tmp: %w", err)
	}

	// Extract to ${tmpDir}
	err = utils.ExtractZip(newPath, tmpDir)
	if err != nil {
		return fmt.Errorf("failed to extract darwin update: %w", err)
	}

	// Remove the zip
	err = os.Remove(newPath)
	if err != nil {
		return fmt.Errorf("failed to remove darwin update zip: %w", err)
	}

	// The extracted .app dir should be ${tmpDir}/${config.AppName}.app
	extractedAppDir := filepath.Join(tmpDir, fmt.Sprintf("%s.app", a.config.AppName))

	// Move the extracted .app dir to .basename.new
	err = os.Rename(extractedAppDir, newPath)
	if err != nil {
		return fmt.Errorf("failed to move darwin update: %w", err)
	}

	// Remove the ${tmpDir} dir
	err = os.RemoveAll(fmt.Sprintf("%s.tmp", newPath))
	if err != nil {
		return fmt.Errorf("failed to remove darwin update tmp: %w", err)
	}

	err = selfupdate.CommitBinary(options)
	if err != nil {
		return fmt.Errorf("failed to commit darwin update: %w", err)
	}

	return nil
}

func (a *DarwinAppApply) OnExit(restart bool) error {
	if restart {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		executable, err := exec.LookPath(os.Args[0])
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		_, err = os.StartProcess(executable, os.Args, &os.ProcAttr{
			Dir:   wd,
			Env:   os.Environ(),
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Sys:   &syscall.SysProcAttr{},
		})

		if err != nil {
			return fmt.Errorf("failed to relaunch after update: %w", err)
		}
	}
	return nil
}
