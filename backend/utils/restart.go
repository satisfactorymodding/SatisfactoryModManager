package utils

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Restart() error {
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
		return fmt.Errorf("failed to start process: %w", err)
	}
	return nil
}
