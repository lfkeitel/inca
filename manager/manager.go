package manager

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

// A Program represents an external executable that will be managed and monitored
type Program struct {
	ConnType          string
	Path              string
	Exec              string
	Exit              chan bool
	AttemptRestarts   int
	attemptedRestarts int
	status            int
	cmd               *exec.Cmd
	conn              net.Conn
}

// Start initially starts the Poller and makes any initial adjustments to the object.
// Start returns nil on a successful start.
func (p *Program) Start() error {
	if p.ConnType != "tcp" && p.ConnType != "unix" {
		return errors.New("Invalid connection type")
	}
	if p.Exit == nil {
		p.Exit = make(chan bool)
	}
	return p.start()
}

// Restart will stop the currently running process and start it again
func (p *Program) Restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	return p.start()
}

// Underlying function that starts the command and dials to the socket or IP port
func (p *Program) start() error {
	if _, err := os.Stat(p.Exec); os.IsNotExist(err) {
		return errors.New("Executable '" + p.Exec + "' not found")
	}

	p.status = 1
	p.cmd = exec.Command(p.Exec, p.ConnType, p.Path)
	if err := p.cmd.Start(); err != nil {
		p.status = 2
		return err
	}

	// Allow time for the poller to start listening
	time.Sleep(2 * time.Second)

	c, err := net.Dial(p.ConnType, p.Path)
	if err != nil {
		p.status = 2
		return err
	}
	p.conn = c
	p.status = 0
	go p.monitorPoller()
	return nil
}

// Stop will stop the currently running poller
func (p *Program) Stop() error {
	if p.cmd != nil {
		// Temporarily set AttemptRestarts to 0
		// so it doesn't trigger and automatic restart by
		// the monitor
		restarts := p.AttemptRestarts
		p.AttemptRestarts = 0
		defer func() { p.AttemptRestarts = restarts }()

		p.cmd.Process.Signal(os.Interrupt)
		select {
		case <-p.Exit:
			break
		case <-time.After(3 * time.Second):
			if err := p.cmd.Process.Kill(); err != nil {
				return err
			}
			break
		}
		p.status = 2
	}

	return nil
}

// monitorPoller will sit and wait for the process to exit. If it exits with an error code,
// it will output it to the console. If the poller has AttemptRestart set to true, monitorPoller
// will attempt a restart
func (p *Program) monitorPoller() {
	err := p.cmd.Wait()
	if err != nil {
		fmt.Println(err.Error())
	}
	p.Exit <- true
	p.cmd = nil
	p.status = 2

	if p.AttemptRestarts > 0 && p.AttemptRestarts < p.attemptedRestarts {
		if err = p.Restart(); err != nil {
			p.attemptedRestarts++
		}
	}
	return
}

func (p *Program) Read() []byte {
	buf := make([]byte, 1024)
	n, err := p.conn.Read(buf[:])
	if err != nil {
		return nil
	}
	return buf[0:n]
}

// Conn returns the Underlying network connection for Program
func (p *Program) Conn() net.Conn {
	return p.conn
}

func (p *Program) Write(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

// Close will close the connection with the process and stop it
func (p *Program) Close() {
	p.conn.Close()
	p.Stop()
	return
}

// Status returns the current status of the Program
func (p *Program) Status() int {
	return p.status
}
