package configs

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"sync"
	"syscall"

	"github.com/dragonrider23/inca/common"
)

var (
	pollerMutex sync.Mutex
	jobIndex    = make(map[int]chan *common.PollerResponse)
	nextJobID   = 0
)

func readPollerResponses() {
	c := poller.Conn()
	for {
		dec := gob.NewDecoder(c)
		var r common.PollerResponse
		err := dec.Decode(&r)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Decode error: " + err.Error())
			}
			continue
		}
		if s, ok := jobIndex[r.ID]; ok {
			ch := r.PrepareReceived()
			go clearIndex(r.ID, ch)
			s <- &r
		}
	}
}

func clearIndex(i int, c <-chan bool) {
	<-c
	delete(jobIndex, i)
}

func sendJobToPoller(cmd common.PollerJob) (<-chan *common.PollerResponse, error) {
	pollerMutex.Lock()
	defer pollerMutex.Unlock()

	cmd.ID = getNextJobID()
	jobIndex[cmd.ID] = make(chan *common.PollerResponse, 1)

	var d bytes.Buffer
	enc := gob.NewEncoder(&d)
	if err := enc.Encode(cmd); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if err := poller.Write(d.Bytes()); err != nil {
		if err.(*net.OpError).Err == syscall.EPIPE {
			return nil, errors.New("Poller not listening")
		}
		return nil, err
	}

	return jobIndex[cmd.ID], nil
}

func getNextJobID() int {
	nextJobID++
	return nextJobID
}
