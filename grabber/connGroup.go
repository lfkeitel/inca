package grabber

import (
	"github.com/dragonrider23/infrastructure-config-archive/common"
)

type connGroup struct {
	numOfConnections int
	goChan           chan bool
	conf             common.Config
}

func newConnGroup(conf common.Config) connGroup {
	return connGroup{
		conf: conf,
	}
}

func (c *connGroup) add(delta int) {
	if c.goChan == nil {
		c.goChan = make(chan bool)
	}
	c.numOfConnections += delta
	return
}

func (c *connGroup) done() {
	c.add(-1)
	finishedDevices++
	if c.numOfConnections < c.conf.MaxSimultaneousConn {
		c.goChan <- true
	}
	return
}

func (c *connGroup) wait() {
	if c.numOfConnections < c.conf.MaxSimultaneousConn {
		return
	}
	<-c.goChan
	return
}
