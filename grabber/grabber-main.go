package grabber

import (
	"fmt"
	"time"

	"github.com/dragonrider23/go-logger"
	"github.com/dragonrider23/infrastructure-config-archive/comm"
	"github.com/dragonrider23/infrastructure-config-archive/targz"
)

var appLogger *logger.Logger
var stdOutLogger *logger.Logger
var configGrabRunning bool
var conf comm.Config

var totalDevices = 0
var finishedDevices = 0

func init() {
	appLogger = logger.New("grabber").Verbose(3).Path("logs/main/")
	stdOutLogger = logger.New("execStdOut").Path("logs/main/")
	configGrabRunning = false
}

func LoadConfig(config comm.Config) {
	conf = config
	return
}

func PerformConfigGrab() {
	if configGrabRunning {
		appLogger.Error("Job already running")
		return
	}

	startTime := time.Now()
	configGrabRunning = true
	defer func() { configGrabRunning = false }()

	// Clean up tftp directory
	removeDir(conf.FullConfDir)

	hosts, err := loadDeviceList(conf)
	if err != nil {
		appLogger.Error(err.Error())
		return
	}

	dtypes, err := loadDeviceTypes(conf)
	if err != nil {
		appLogger.Error(err.Error())
		return
	}

	totalDevices = len(hosts)
	finishedDevices = 0
	dateSuffix := time.Now().Format("2006012")

	grabConfigs(hosts, dtypes, dateSuffix, conf)
	tarGz.TarGz("archive/"+dateSuffix+".tar.gz", conf.FullConfDir)

	endTime := time.Now()
	logText := fmt.Sprintf("Config grab took %s", endTime.Sub(startTime).String())
	appLogger.Info(logText)
	comm.UserLogInfo(logText)
	return
}

func PerformSingleRun(name, hostname, brand, method string) {
	if configGrabRunning {
		appLogger.Error("Job already running")
		return
	}

	startTime := time.Now()
	configGrabRunning = true
	defer func() { configGrabRunning = false }()

	hosts := make([]host, 1)

	hosts[0] = host{
		name:    name,
		address: hostname,
		dtype:   brand,
		method:  method,
	}

	dtypes, err := loadDeviceTypes(conf)
	if err != nil {
		appLogger.Error(err.Error())
		return
	}

	totalDevices = 1
	finishedDevices = 0
	dateSuffix := time.Now().Format("2006012")

	grabConfigs(hosts, dtypes, dateSuffix, conf)
	tarGz.TarGz("archive/"+dateSuffix+".tar.gz", conf.FullConfDir)

	endTime := time.Now()
	logText := fmt.Sprintf("Config grab took %s", endTime.Sub(startTime).String())
	appLogger.Info(logText)
	comm.UserLogInfo(logText)
	return
}

func IsRunning() bool {
	return configGrabRunning
}

func Remaining() (total, finished int) {
	if !configGrabRunning {
		if totalDevices == 0 {
			hosts, err := loadDeviceList(conf)
			if err != nil {
				appLogger.Error(err.Error())
				return
			}
			totalDevices = len(hosts)
		}

		if finishedDevices == 0 {
			finishedDevices = -1
		}
	}

	return totalDevices, finishedDevices
}
