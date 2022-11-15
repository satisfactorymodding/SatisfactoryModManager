package apply

import (
	"io"
	"os"
	"syscall"

	"github.com/minio/selfupdate"
	"github.com/pkg/errors"
)

type SingleFileApply struct {
}

func MakeSingleFileApply() *SingleFileApply {
	return &SingleFileApply{}
}

func (a *SingleFileApply) Apply(file io.Reader) error {
	err := selfupdate.Apply(file, selfupdate.Options{})
	if err != nil {
		return errors.Wrap(err, "failed to apply singlefile update")
	}
	return nil
}

func (a *SingleFileApply) OnExit(restart bool) error {
	if restart {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		executable, err := os.Executable()
		if err != nil {
			return errors.Wrap(err, "failed to get executable path")
		}

		_, err = os.StartProcess(executable, os.Args, &os.ProcAttr{
			Dir:   wd,
			Env:   os.Environ(),
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Sys:   &syscall.SysProcAttr{},
		})

		if err != nil {
			return errors.Wrap(err, "failed to relaunch after update")
		}
	}
	return nil
}
