package configs

import (
	"fmt"
	"time"

	"github.com/dragonrider23/go-logger"

	"github.com/dragonrider23/infrastructure-config-archive/common"
	"github.com/dragonrider23/infrastructure-config-archive/devices"
)

var appLogger *logger.Logger
var stdOutLogger *logger.Logger
var configGrabRunning bool
var conf common.Config

var totalDevices = 0
var finishedDevices = 0

func init() {
	appLogger = logger.New("grabber").Verbose(3).Path("logs/main/")
	stdOutLogger = logger.New("execStdOut").Path("logs/main/")
	configGrabRunning = false
}

func LoadConfig(config common.Config) {
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

	hosts, err := devices.GetDevicesForConfigs()
	if err != nil {
		appLogger.Error(err.Error())
		return
	}
	//fmt.Printf("%#v", hosts)

	connProfiles, err := devices.GetConnProfiles()
	if err != nil {
		appLogger.Error(err.Error())
		return
	}
	//fmt.Printf("%#v", connProfiles)

	totalDevices = len(hosts)
	finishedDevices = 0

	grabConfigs(hosts, connProfiles)

	endTime := time.Now()
	logText := fmt.Sprintf("Config grab took %s", endTime.Sub(startTime).String())
	appLogger.Info(logText)
	// common.UserLogInfo(logText)
	return
}

func IsRunning() bool {
	return configGrabRunning
}

func Remaining() (total, finished int) {
	if !configGrabRunning {
		if totalDevices == 0 {
			hosts, err := devices.GetAllDevices()
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
