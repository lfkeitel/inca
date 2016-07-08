package comm

import (
	"os"
	"time"

	"github.com/lfkeitel/go-logger"
)

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
	logger.New("endUserLog").NoStdout().Raw().Path("logs/endUser/")
}

func ReverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func UserLogInfo(format string, v ...interface{}) {
	format = "INFO:-:" + time.Now().Format("2006-01-02 15:04:05") + ":-:" + format
	logger.Get("endUserLog").Log("log", logger.Cyan, format, v...)
	return
}

func UserLogWarning(format string, v ...interface{}) {
	format = "WARNING:-:" + time.Now().Format("2006-01-02 15:04:05") + ":-:" + format
	logger.Get("endUserLog").Log("log", logger.Cyan, format, v...)
	return
}

func UserLogError(format string, v ...interface{}) {
	format = "ERROR:-:" + time.Now().Format("2006-01-02 15:04:05") + ":-:" + format
	logger.Get("endUserLog").Log("log", logger.Cyan, format, v...)
	return
}

func UserLogFatal(format string, v ...interface{}) {
	format = "FATAL:-:" + time.Now().Format("2006-01-02 15:04:05") + ":-:" + format
	logger.Get("endUserLog").Log("log", logger.Cyan, format, v...)
	os.Exit(1)
	return
}
