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

func (a *SingleFileApply) Download(file io.Reader, checksum []byte) error {
	err := selfupdate.PrepareAndCheckBinary(file, selfupdate.Options{
		Checksum: checksum,
	})
	if err != nil {
		return fmt.Errorf("failed to download singlefile update: %w", err)
	}
	return nil
}

func (a *SingleFileApply) Apply(restart bool) error {
	err := selfupdate.CommitBinary(selfupdate.Options{})
	if err != nil {
		return fmt.Errorf("failed to commit singlefile update: %w", err)
	}
	if restart {
		err := utils.Restart()
		if err != nil {
			return fmt.Errorf("failed to relaunch after update: %w", err)
		}
	}
	return nil
}
