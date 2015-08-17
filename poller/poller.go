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
	Data interface{}
}

// Response represents a response from the poller for a command
type Response struct {
	Error string
	Data  interface{}
}

var (
	acceptNew = true
	wg        = sync.WaitGroup{}
)

// Process receives job j and starts a goroutine to process a response. The response will be sent
// on the returned channel. An error will return if the job didn't start successfully
func Process(j Job) (<-chan *Response, error) {
	if !acceptNew {
		return nil, errors.New("Not accepting new jobs")
	}
	out := make(chan *Response, 1)

	switch j.Cmd {
	case "echo":
		wg.Add(1)
		go func() {
			defer wg.Done()
			echoJob(j, out)
		}()
	case "sleep":
		wg.Add(1)
		go func() {
			defer wg.Done()
			sleepJob(j, out)
		}()
	default:
		return nil, errors.New("Bad command")
	}
	return out, nil
}

// Close stops the poller from accepting new jobs and waits for the currently running ones
// to exit
func Close() {
	acceptNew = false
	fmt.Println("Waiting for jobs to finish")
	wg.Wait()
}

func echoJob(j Job, out chan<- *Response) {
	out <- &Response{
		Error: "",
		Data:  j.Data,
	}
}

func sleepJob(j Job, out chan<- *Response) {
	<-time.After(time.Duration(j.Data.(int)) * time.Second)
	out <- &Response{
		Error: "",
		Data:  "",
	}
}
