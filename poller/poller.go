package poller

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Job represents a command given to the poller
type Job struct {
	Cmd  string
	Data map[string]interface{}
}

// Response represents a response from the poller for a command
type Response struct {
	Error string
	Data  interface{}
}

var (
	acceptNew       = true
	wg              = sync.WaitGroup{}
	errBadJob       = errors.New("Bad job")
	errNotAccepting = errors.New("Not accepting new jobs")
)

// Process receives job j and starts a goroutine to process a response. The response will be sent
// on the returned channel. An error will return if the job didn't start successfully
func Process(j Job) (<-chan *Response, <-chan error, error) {
	if !acceptNew {
		return nil, nil, errNotAccepting
	}
	out := make(chan *Response, 1)
	err := make(chan error, 1)

	switch j.Cmd {
	case "echo":
		wg.Add(1)
		go func() {
			defer wg.Done()
			echoJob(j, out, err)
		}()
	case "sleep":
		wg.Add(1)
		go func() {
			defer wg.Done()
			sleepJob(j, out, err)
		}()
	case "poll":
		wg.Add(1)
		go func() {
			defer wg.Done()
			configJob(j, out, err)
		}()
	default:
		return nil, nil, errBadJob
	}
	return out, err, nil
}

// Close stops the poller from accepting new jobs and waits for the currently running ones
// to exit
func Close() {
	acceptNew = false
	fmt.Println("Waiting for jobs to finish")
	wg.Wait()
}

func echoJob(j Job, out chan<- *Response, e chan<- error) {
	out <- &Response{
		Error: "",
		Data:  j.Data["echo"],
	}
}

func sleepJob(j Job, out chan<- *Response, e chan<- error) {
	<-time.After(time.Duration(j.Data["timeout"].(int)) * time.Second)
	out <- &Response{
		Error: "",
		Data:  "",
	}
}
