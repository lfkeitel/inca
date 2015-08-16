package configs

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dragonrider23/inca/common"
	"github.com/dragonrider23/inca/manager"
)

var pollerRWMutex sync.RWMutex

func startPoller() error {
	var path string
	config.Poller.Connection = strings.ToLower(config.Poller.Connection)

	switch config.Poller.Connection {
	case "tcp":
		if config.Poller.Path != "" && config.Poller.Port != "" {
			path = config.Poller.Path + ":" + config.Poller.Port
		} else {
			return errors.New("Path and port must be defined for IP connection")
		}
		break
	case "unix":
		if config.Poller.Path == "" {
			path = filepath.Join(os.TempDir(), "inca-socket.sock")
		}
		break
	default:
		return errors.New("Poller Connection must be either ip or unix")
	}

	poller = &manager.Program{
		ConnType:        config.Poller.Connection,
		Path:            path,
		Exec:            "./poller",
		AttemptRestarts: 3,
	}
	if err := poller.Start(); err != nil {
		return err
	}
	return nil
}

// PollerStatus returns an int depending on the current status of the poller
// 0 = Running, OK
// 1 = Starting
// 2 = Stopped
func PollerStatus() int {
	return poller.Status()
}

// HeartBeat sends a heartbeat to the poller and returns the poller's response
func HeartBeat() string {
	cmd := common.PollerCommand{
		Cmd:  "echo",
		Data: "heartbeat",
	}

	err := sendCommandToPoller(cmd)
	if err != nil {
		return "There was an error"
	}

	r, err := readFromPoller()
	if err != nil {
		return "There was an error"
	}
	return r.Data.(string)
}

func readFromPoller() (*common.PollerResponse, error) {
	pollerRWMutex.RLock()
	c := poller.Conn()
	dec := gob.NewDecoder(c)

	var r common.PollerResponse
	err := dec.Decode(&r)
	if err != nil {
		return nil, err
	}
	pollerRWMutex.RUnlock()
	return &r, nil
}

func sendCommandToPoller(cmd common.PollerCommand) error {
	pollerRWMutex.Lock()
	var d bytes.Buffer
	enc := gob.NewEncoder(&d)
	if err := enc.Encode(cmd); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := poller.Write(d.Bytes()); err != nil {
		fmt.Println(err.Error())
		return err
	}
	pollerRWMutex.Unlock()

	return nil
}
