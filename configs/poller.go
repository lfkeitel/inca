package configs

import (
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/dragonrider23/inca/common"
	"github.com/dragonrider23/inca/manager"
)

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

	if config.Poller.Exec == "" {
		config.Poller.Exec = "./poller"
	}

	poller = &manager.Program{
		ConnType:        config.Poller.Connection,
		Path:            path,
		Exec:            config.Poller.Exec,
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
func HeartBeat() (string, error) {
	cmd := common.PollerJob{
		Cmd:  "echo",
		Data: "heartbeat",
	}

	c, err := sendJobToPoller(cmd)
	if err != nil {
		return "", err
	}

	r := <-c
	r.Received()
	return r.Data.(string), nil
}

func pollerSendConfig() error {
	gob.Register(config)
	cmd := common.PollerJob{
		Cmd:  "config",
		Data: config,
	}

	c, err := sendJobToPoller(cmd)
	if err != nil {
		return err
	}

	r := <-c
	r.Received()
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}
