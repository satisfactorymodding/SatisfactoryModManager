package apply

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/minio/selfupdate"
)

type SingleFileApply struct{}

func MakeSingleFileApply() *SingleFileApply {
	return &SingleFileApply{}
}

func (a *SingleFileApply) Apply(file io.Reader, checksum []byte) error {
	err := selfupdate.Apply(file, selfupdate.Options{
		Checksum: checksum,
	})
	if err != nil {
		return fmt.Errorf("failed to apply singlefile update: %w", err)
	}
	return nil
}

func (a *SingleFileApply) OnExit(restart bool) error {
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
