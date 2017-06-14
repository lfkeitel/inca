package common

import "github.com/lfkeitel/verbose"

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

func init() {
	userLog := verbose.New("endUserLog")

	fileLogger, err := verbose.NewFileHandler("logs/endUser/")
	if err != nil {
		panic("Failed to open logging directory")
	}

	userLog.AddHandler("file", fileLogger)
}

func ReverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func UserLogInfo(format string, v ...interface{}) {
	verbose.Get("endUserLog").Infof(format, v...)
}

func UserLogWarning(format string, v ...interface{}) {
	verbose.Get("endUserLog").Warningf(format, v...)
}

func UserLogError(format string, v ...interface{}) {
	verbose.Get("endUserLog").Errorf(format, v...)
}

func UserLogFatal(format string, v ...interface{}) {
	verbose.Get("endUserLog").Fatalf(format, v...)
}
