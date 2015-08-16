package configs

import (
	"fmt"

	"github.com/dragonrider23/inca/common"
	"github.com/dragonrider23/inca/manager"
)

var (
	config common.Config
	poller *manager.Program
)

// Prepare sets the configuration for the package and starts the poller executable
func Prepare(conf common.Config) error {
	config = conf
	return startPoller()
}

// Stop shuts down the poller
func Stop() {
	fmt.Println("Stopping poller")
	poller.Stop()
}
