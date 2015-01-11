package interfaces

type Config struct {
    Tftphost       string
    RemoteUsername string
    RemotePassword string
    EnablePassword string
    DeviceListFile string
    FullConfDir    string
    MaxSimultaneousConn int
    Server serverConf
}

type serverConf struct {
    BindAddress string
    BindPort    int
    BaseDir     string
}
