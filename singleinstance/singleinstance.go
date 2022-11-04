//go:build !bindings

package singleinstance

import (
	"encoding/json"
	"os"
	"path"

	"bitbucket.org/avd/go-ipc/mq"
	"github.com/juju/fslock"
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
			panic(err)
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
	mq.Destroy("SatisfactoryModManager")
	messageQueue, err := mq.New("SatisfactoryModManager", os.O_CREATE|mq.O_NONBLOCK, 0666)
	if err != nil {
		panic(err)
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
			if !mq.IsTemporary(err) {
				panic(err)
			}
		}
		if l > 0 {
			trimmedData := data[:l]
			var args []string
			err := json.Unmarshal(trimmedData, &args)
			if err != nil {
				panic(err)
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
		return err
	}
	defer messageQueue.Close()
	argsBytes, err := json.Marshal(os.Args)
	if err != nil {
		return err
	}
	err = messageQueue.Send(argsBytes)
	if err != nil {
		return err
	}
	return nil
}
