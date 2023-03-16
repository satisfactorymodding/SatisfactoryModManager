//go:build !bindings

package singleinstance

import (
	"encoding/json"
	"os"
	"path"

	"bitbucket.org/avd/go-ipc/mq"
	"github.com/juju/fslock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Similar behaviour to electron's requestSingleInstanceLock.
// Checks if there is already an instance running and if so, passes the arguments to it.
func RequestSingleInstanceLock() bool {
	lockFile := getLockFile()
	lock := fslock.New(lockFile)
	err := lock.TryLock()
	if err != nil {
		err = sendArgs()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to send arguments to first instance")
		}
		return false
	}
	return true
}

var listening = false

func ListenForSecondInstance() {
	if listening {
		return
	}
	_ = mq.Destroy("SatisfactoryModManager")
	messageQueue, err := mq.New("SatisfactoryModManager", os.O_CREATE|mq.O_NONBLOCK, 0o666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create message queue")
	}
	data := make([]byte, 512)
	listening = true
	for {
		l, err := messageQueue.Receive(data)
		if err != nil {
			if err.Error() == "the message is too long" {
				data = make([]byte, 2*len(data)) // double the size hoping to fit the message
				continue
			}
			if !mq.IsTemporary(err) && err.Error() != "MSGRCV: no message of desired type" {
				log.Error().Err(err).Msg("Failed to receive message")
			}
		}
		if l > 0 {
			trimmedData := data[:l]
			var args []string
			err := json.Unmarshal(trimmedData, &args)
			if err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal arguments")
			}
			OnSecondInstance(args)
		}
	}
}

func getLockFile() string {
	return path.Join(os.TempDir(), "SatisfactoryModManager.lock")
}

func sendArgs() error {
	messageQueue, err := mq.Open("SatisfactoryModManager", mq.O_NONBLOCK)
	if err != nil {
		return errors.Wrap(err, "Failed to open message queue")
	}
	defer messageQueue.Close()
	argsBytes, err := json.Marshal(os.Args)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal arguments")
	}
	err = messageQueue.Send(argsBytes)
	if err != nil {
		return errors.Wrap(err, "Failed to send arguments")
	}
	return nil
}
