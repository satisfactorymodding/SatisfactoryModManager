package apply

import (
	"crypto/sha256"
	"io"
	"os"
	"os/exec"
	"slices"

	"github.com/pkg/errors"
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

func (a *NsisApply) Apply(file io.Reader, checksum []byte) error {
	err := a.writeInstaller(file)
	if err != nil {
		return err
	}

	return a.checkHash(checksum)
}

func (a *NsisApply) writeInstaller(file io.Reader) error {
	f, err := os.OpenFile(a.config.InstallerDownloadPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return errors.Wrap(err, "failed to open installer file")
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return errors.Wrap(err, "failed to write installer file")
	}
	return nil
}

func (a *NsisApply) checkHash(checksum []byte) error {
	if checksum == nil {
		return nil
	}
	f, err := os.Open(a.config.InstallerDownloadPath)
	if err != nil {
		return errors.Wrap(err, "failed to open installer file")
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "failed to read installer file")
	}
	installerHash := sha256.Sum256(data)
	if !slices.Equal(installerHash[:], checksum) {
		return errors.Errorf("installer hash does not match")
	}
	return nil
}

func (a *NsisApply) OnExit(restart bool) error {
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
		return errors.Wrap(err, "failed to start nsis installer")
	}
	err = cmd.Process.Release()
	if err != nil {
		return errors.Wrap(err, "failed to release nsis installer process")
	}
	return nil
}
