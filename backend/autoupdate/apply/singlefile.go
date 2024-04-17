package apply

import (
	"fmt"
	"io"

	"github.com/minio/selfupdate"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
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
		err := utils.Restart()
		if err != nil {
			return fmt.Errorf("failed to relaunch after update: %w", err)
		}
	}
	return nil
}
