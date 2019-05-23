package server

import (
	"encoding/json"
	"fmt"
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
	Configs      []string `json:"configs"`
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

	http.HandleFunc("/", spaServer)
	http.HandleFunc("/api/", apiHandler)

	err := http.ListenAndServe(conf.Server.BindAddress+":"+strconv.Itoa(conf.Server.BindPort), nil)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
}

func spaServer(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	path := filepath.Join("frontend", upath)
	if !common.FileExists(path) {
		// Redirect all requests to the index page
		path = filepath.Join("frontend", "index.html")
	}
	http.ServeFile(w, r, path)
}

// Get a list of all devices in the config.FullConfDir directory
func getDeviceList() deviceList {
	type host struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Dtype   string `json:"dtype"`
		Method  string `json:"method"`
	}

	configFileList, _ := ioutil.ReadDir(config.Paths.ConfDir)

	deviceConfigs := deviceList{}

	for _, file := range configFileList {
		if !file.IsDir() {
			continue
		}

		hostdir := filepath.Join(config.Paths.ConfDir, file.Name())
		filename := filepath.Join(hostdir, "_metadata.json")
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			continue
		}

		var device host
		if err := json.Unmarshal(data, &device); err != nil {
			continue
		}

		dircontents, _ := ioutil.ReadDir(hostdir)
		files := make([]string, 0, len(dircontents)-1)
		for _, f := range dircontents {
			if f.Name() != "_metadata.json" {
				files = append(files, f.Name())
			}
		}

		deviceConfigs.Devices = append(deviceConfigs.Devices, deviceConfigFile{
			Path:         fmt.Sprintf("%s-%s", device.Name, device.Address),
			Name:         device.Name,
			Address:      device.Address,
			Proto:        device.Dtype,
			Manufacturer: device.Method,
			Configs:      files,
		})
	}

	return deviceConfigs
}
