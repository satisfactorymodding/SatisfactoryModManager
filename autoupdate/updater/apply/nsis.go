package apply

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type NsisApply struct {
	config NsisApplyConfig
}

type NsisApplyConfig struct {
	InstallerDownloadPath string
	Elevation             bool
}

func MakeNsisApply(config NsisApplyConfig) *NsisApply {
	return &NsisApply{
		config: config,
	}
}

func (a *NsisApply) Apply(file io.Reader) error {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read installer file")
	}
	err = os.WriteFile(a.config.InstallerDownloadPath, fileContent, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to write installer file")
	}
	return nil
}

func (a *NsisApply) OnExit(restart bool) error {
	arguments := []string{"/S"}
	if restart {
		// TODO implement this in my installer
		arguments = append(arguments, "--force-run")
	}
	cmd := exec.Command(a.config.InstallerDownloadPath, arguments...)
	err := cmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start nsis installer")
	}
	err = cmd.Process.Release()
	if err != nil {
		return errors.Wrap(err, "failed to release nsis installer process")
	}
	return nil
}
