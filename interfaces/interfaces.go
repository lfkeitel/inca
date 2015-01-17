package interfaces

type Config struct {
	RemoteUsername      string
	RemotePassword      string
	EnablePassword      string
	DeviceListFile      string
	DeviceTypeFile      string
	FullConfDir         string
	MaxSimultaneousConn int
	Server              serverConf
}

type serverConf struct {
	BindAddress string
	BindPort    int
}
