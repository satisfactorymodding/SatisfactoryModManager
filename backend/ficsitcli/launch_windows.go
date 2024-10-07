package ficsitcli

import (
	"os/exec"
	"syscall"
)

func (f *ficsitCLI) executeLaunchCommand(launchPath []string) ([]byte, string, error) {
	cmd := exec.Command(launchPath[0], launchPath[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	return out, cmd.String(), err
}
