package grabber

type maxChanGroup struct {
	numOfConnections int
	goChan           chan bool
	maxConnections   int
}

func newMaxChanGroup(max int) maxChanGroup {
	return maxChanGroup{
		maxConnections: max,
	}
}

func (c *maxChanGroup) add(delta int) {
	if c.goChan == nil {
		c.goChan = make(chan bool)
	}
	c.numOfConnections += delta
	return
}

func (c *maxChanGroup) done() {
	c.add(-1)
	finishedDevices++
	if c.numOfConnections < c.maxConnections {
		c.goChan <- true
	}
	return
}

func (c *maxChanGroup) wait() {
	if c.numOfConnections < c.maxConnections {
		return
	}
	<-c.goChan
	return
}
