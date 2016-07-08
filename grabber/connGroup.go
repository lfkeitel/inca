package grabber

import (
	"github.com/lfkeitel/inca/comm"
)

type connGroup struct {
	numOfConnections int
	goChan           chan bool
	conf             comm.Config
}

func newConnGroup(conf comm.Config) connGroup {
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
