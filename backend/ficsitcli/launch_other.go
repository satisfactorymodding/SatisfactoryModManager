//go:build !windows

package ficsitcli

import (
	"os/exec"
)

func (f *ficsitCLI) executeLaunchCommand(launchPath []string) ([]byte, string, error) {
	cmd := exec.Command(launchPath[0], launchPath[1:]...)
	out, err := cmd.CombinedOutput()
	return out, cmd.String(), err
}
