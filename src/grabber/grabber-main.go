package grabber

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/lfkeitel/inca/src/common"
	tarGz "github.com/lfkeitel/inca/src/targz"
	"github.com/lfkeitel/verbose"
)

var appLogger *verbose.Logger
var stdOutLogger *verbose.Logger
var configGrabRunning = false
var conf common.Config

var totalDevices = 0
var finishedDevices = 0
var stage = ""

func init() {
	appLogger = verbose.New("grabber")
	stdOutLogger = verbose.New("execStdOut")

	fileLogger, err := verbose.NewFileHandler("logs/main.log")
	if err != nil {
		panic("Failed to open logging directory")
	}
	fileLogger.SetMinLevel(verbose.LogLevelInfo)

	appLogger.AddHandler("file", fileLogger)
	stdOutLogger.AddHandler("file", fileLogger)
}

func LoadConfig(config common.Config) { conf = config }

func PerformConfigGrab() {
	if configGrabRunning {
		appLogger.Error("Job already running")
		return
	}

	startTime := time.Now()
	configGrabRunning = true
	defer func() {
		configGrabRunning = false
		stage = ""
	}()

	if conf.Hooks.PreScript != "" {
		appLogger.Info("Running pre script")
		stage = "pre-script"
		if err := exec.Command(conf.Hooks.PreScript).Run(); err != nil {
			appLogger.Error(err)
		}
	}

	stage = "loading-configuration"
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

	existing, err := loadCurrentConfigs(conf)
	if err != nil {
		appLogger.Error(err.Error())
		return
	}

	totalDevices = len(hosts)
	finishedDevices = 0
	dateSuffix := time.Now().Format("20060102")

	stage = "grabbing"
	grabConfigs(hosts, dtypes, dateSuffix, conf, existing)
	tarGz.TarGz("archive/"+dateSuffix+".tar.gz", conf.FullConfDir)

	if conf.Hooks.PostScript != "" {
		appLogger.Info("Running post script")
		stage = "post-script"
		if err := exec.Command(conf.Hooks.PostScript).Run(); err != nil {
			appLogger.Error(err)
		}
	}

	endTime := time.Now()
	logText := fmt.Sprintf("Config grab took %s", endTime.Sub(startTime).String())
	appLogger.Info(logText)
	common.UserLogInfo(logText)
	return
}

func PerformSingleRun(name, hostname, brand, method string) {
	if configGrabRunning {
		appLogger.Error("Job already running")
		return
	}

	startTime := time.Now()
	configGrabRunning = true
	defer func() {
		configGrabRunning = false
		stage = ""
	}()
	name = strings.Replace(name, "-", "_", -1)

	if conf.Hooks.PreScript != "" {
		appLogger.Info("Running pre script")
		stage = "pre-script"
		if err := exec.Command(conf.Hooks.PreScript).Run(); err != nil {
			appLogger.Error(err)
		}
	}

	stage = "loading-configuration"
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

	existing, err := loadCurrentConfigs(conf)
	if err != nil {
		appLogger.Error(err.Error())
		return
	}

	totalDevices = 1
	finishedDevices = 0
	dateSuffix := time.Now().Format("20060102")

	stage = "grabbing"
	grabConfigs(hosts, dtypes, dateSuffix, conf, existing)
	tarGz.TarGz("archive/"+dateSuffix+".tar.gz", conf.FullConfDir)

	if conf.Hooks.PostScript != "" {
		appLogger.Info("Running post script")
		stage = "post-script"
		if err := exec.Command(conf.Hooks.PostScript).Run(); err != nil {
			appLogger.Error(err)
		}
	}

	endTime := time.Now()
	logText := fmt.Sprintf("Config grab took %s", endTime.Sub(startTime).String())
	appLogger.Info(logText)
	common.UserLogInfo(logText)
	return
}

func IsRunning() bool {
	return configGrabRunning
}

func remainingDeviceCount() (total, finished int) {
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

type State struct {
	Running         bool
	Total, Finished int
	Stage           string
}

func CurrentState() State {
	total, finished := remainingDeviceCount()

	return State{
		Running:  configGrabRunning,
		Total:    total,
		Finished: finished,
		Stage:    stage,
	}
}
