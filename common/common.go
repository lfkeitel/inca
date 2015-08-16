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
	Connection string
	Path       string
	Port       string
}

// PollerCommand represents a command given to the poller
type PollerCommand struct {
	Cmd  string
	Data interface{}
}

// PollerResponse represents a response from the poller for a command
type PollerResponse struct {
	Error string
	Data  interface{}
}
