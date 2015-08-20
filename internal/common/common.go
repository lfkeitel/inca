package common

// A Configuration represents the structure of the configuration file
type Configuration struct {
	MaxSimultaneousConn int
	DataDir             string
	DashboardLogLevel   string
	Server              serverConf
	Database            databaseConf
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
