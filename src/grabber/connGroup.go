package grabber

type connGroup struct {
	numOfConnections int
	goChan           chan bool
	max              int
}

func newConnGroup(max int) connGroup {
	return connGroup{
		max: max,
	}
}

func (c *connGroup) add(delta int) {
	if c.goChan == nil {
		c.goChan = make(chan bool)
	}
	c.numOfConnections += delta
}

func (c *connGroup) done() {
	c.add(-1)
	if c.numOfConnections < c.max {
		c.goChan <- true
	}
}

func (c *connGroup) wait() {
	if c.numOfConnections < c.max {
		return
	}
	<-c.goChan
}
