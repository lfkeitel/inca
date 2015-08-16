package common

import ()

// A Config represents the structure of the configuration file
type Config struct {
	MaxSimultaneousConn int
	DataDir             string
	DashboardLogLevel   string
	Server              serverConf
	Database            databaseConf
	Poller              pollerConf
}

type serverConf struct {
	BindAddress string
	BindPort    int
}

type databaseConf struct {
	Address  string
	Port     int
	Username string
	Password string
	Name     string
}

type pollerConf struct {
	Exec       string
	Connection string
	Path       string
	Port       string
}

// PollerJob represents a command given to the poller
type PollerJob struct {
	ID   int
	Cmd  string
	Data interface{}
}

// PollerResponse represents a response from the poller for a command
type PollerResponse struct {
	ID    int
	Error string
	Data  interface{}
	rec   chan bool
}

// PrepareReceived makes a channel on the response object that's used to notify
// when the receiver has the object
func (p *PollerResponse) PrepareReceived() chan bool {
	p.rec = make(chan bool, 1)
	return p.rec
}

// Received is called by the receiver to notify it has the object
func (p *PollerResponse) Received() {
	p.rec <- true
}
