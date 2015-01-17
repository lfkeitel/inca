package grabber

import (
	"time"

	logger "github.com/dragonrider23/go-logger"
	"github.com/dragonrider23/infrastructure-config-archive/interfaces"
	"github.com/dragonrider23/infrastructure-config-archive/targz"
)

var appLogger *logger.Logger
var stdOutLogger *logger.Logger
var configGrabRunning bool
var conf interfaces.Config

var totalDevices = 0
var finishedDevices = 0

func init() {
	appLogger = logger.New("grabber").Verbose(3)
	stdOutLogger = logger.New("execStdOut")
	configGrabRunning = false
}

func LoadConfig(config interfaces.Config) {
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

	totalDevices = len(hosts)
	finishedDevices = 0
	dateSuffix := time.Now().Format("2006012")

	grabConfigs(hosts, dateSuffix, conf)
	tarGz.TarGz("archive/"+dateSuffix+".tar.gz", conf.FullConfDir)

	endTime := time.Now()
	appLogger.Info("Config grab took %s", endTime.Sub(startTime).String())
	return
}

func PerformSingleRun(name, hostname, brand, proto string) {
	if configGrabRunning {
		appLogger.Error("Job already running")
		return
	}

	startTime := time.Now()
	configGrabRunning = true
	defer func() { configGrabRunning = false }()

	hosts := make([]host, 1)

	hosts[0] = host{
		name:         name,
		address:      hostname,
		manufacturer: brand,
		proto:        proto,
	}

	totalDevices = 1
	finishedDevices = 0
	dateSuffix := time.Now().Format("2006012")

	grabConfigs(hosts, dateSuffix, conf)
	tarGz.TarGz("archive/"+dateSuffix+".tar.gz", conf.FullConfDir)

	endTime := time.Now()
	appLogger.Info("Config grab took %s", endTime.Sub(startTime).String())
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
