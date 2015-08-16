package configs

import (
	"fmt"

	"github.com/dragonrider23/inca/common"
	"github.com/dragonrider23/inca/manager"
)

var (
	config *common.Config
	poller *manager.Program
)

// Prepare sets the configuration for the package and starts the poller executable
func Prepare(conf *common.Config) error {
	config = conf
	if err := startPoller(); err != nil {
		return err
	}

	go readPollerResponses()
	go detectPollerRestart()
	return nil
}

func detectPollerRestart() {
	for {
		<-poller.Started
		if err := pollerSendConfig(); err != nil {
			Stop()
		}
	}
}

// Stop shuts down the poller
func Stop() {
	fmt.Println("Stopping poller")
	poller.Stop()
}
