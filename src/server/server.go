package server

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lfkeitel/inca/src/common"
	"github.com/lfkeitel/verbose"
)

type deviceConfigFile struct {
	Path         string   `json:"path"`
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Proto        string   `json:"proto"`
	ConfText     []string `json:"conf_text"`
	Manufacturer string   `json:"manufacturer"`
}

type deviceList struct {
	Devices []deviceConfigFile `json:"devices"`
}

var appLogger *verbose.Logger
var config *common.Config

// Initialize HTTP server with app configuration and templates
func initServer(configuration *common.Config) {
	config = configuration

	appLogger = verbose.New("httpServer")

	fileLogger, err := verbose.NewFileHandler(filepath.Join(configuration.Paths.LogDir, "server.log"))
	if err != nil {
		panic("Failed to open logging directory")
	}

	appLogger.AddHandler("file", fileLogger)
	appLogger.AddHandler("stdout", verbose.NewStdoutHandler(true))
}

// Start front-end HTTP server
func StartServer(conf *common.Config) {
	initServer(conf)

	logText := "Starting webserver on port " + conf.Server.BindAddress + ":" + strconv.Itoa(conf.Server.BindPort)
	appLogger.Info(logText)
	common.UserLogInfo(logText)

	http.Handle("/", http.FileServer(http.Dir(filepath.Join("frontend"))))
	http.HandleFunc("/api/", apiHandler)

	err := http.ListenAndServe(conf.Server.BindAddress+":"+strconv.Itoa(conf.Server.BindPort), nil)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
}

// Get a list of all devices in the config.FullConfDir directory
func getDeviceList() deviceList {
	configFileList, _ := ioutil.ReadDir(config.Paths.ConfDir)

	deviceConfigs := deviceList{}

	for _, file := range configFileList {
		filename := file.Name()
		if filename[0] == '.' {
			continue
		}
		splitName := strings.Split(filename, "-")      // [0] = name, [1] = datesuffix, [2] = hostname, [3] = manufacturer
		splitProto := strings.Split(splitName[4], ".") // [0] = protocol, [1] = ".conf"

		device := deviceConfigFile{
			Path:         file.Name(),
			Name:         splitName[0],
			Address:      splitName[2],
			Proto:        splitProto[0],
			Manufacturer: splitName[3],
		}
		deviceConfigs.Devices = append(deviceConfigs.Devices, device)
	}

	return deviceConfigs
}
