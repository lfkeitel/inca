package common

import (
	"os"
	//"time"

	"github.com/dragonrider23/go-logger"
)

type Config struct {
	MaxSimultaneousConn int
	DataDir             string
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

func init() {
	logger.New("endUserLog").NoStdout().Raw().Path("logs/endUser/")
}

func UserLogInfo(format string, v ...interface{}) {
	// format = "INFO:-:" + time.Now().Format("2006-01-02 15:04:05") + ":-:" + format
	// logger.Get("endUserLog").Log("log", logger.Cyan, format, v...)
	return
}

func UserLogWarning(format string, v ...interface{}) {
	// format = "WARNING:-:" + time.Now().Format("2006-01-02 15:04:05") + ":-:" + format
	// logger.Get("endUserLog").Log("log", logger.Cyan, format, v...)
	return
}

func UserLogError(format string, v ...interface{}) {
	// format = "ERROR:-:" + time.Now().Format("2006-01-02 15:04:05") + ":-:" + format
	// logger.Get("endUserLog").Log("log", logger.Cyan, format, v...)
	return
}

func UserLogFatal(format string, v ...interface{}) {
	// format = "FATAL:-:" + time.Now().Format("2006-01-02 15:04:05") + ":-:" + format
	// logger.Get("endUserLog").Log("log", logger.Cyan, format, v...)
	os.Exit(1)
	return
}
