package server

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/lfkeitel/inca/common"
	"github.com/lfkeitel/verbose"
)

type deviceConfigFile struct {
	Path         string
	Name         string
	Address      string
	Proto        string
	ConfText     []string
	Manufacturer string
}

type deviceList struct {
	Devices []deviceConfigFile
}

var templates *template.Template
var appLogger *verbose.Logger
var config common.Config

// Initialize HTTP server with app configuration and templates
func initServer(configuration common.Config) {
	config = configuration
	templates = template.Must(template.ParseGlob("server/templates/*.tmpl"))

	appLogger = verbose.New("httpServer")

	fileLogger, err := verbose.NewFileHandler("logs/server/")
	if err != nil {
		panic("Failed to open logging directory")
	}

	appLogger.AddHandler("file", fileLogger)
}

// Start front-end HTTP server
func StartServer(conf common.Config) {
	initServer(conf)

	logText := "Starting webserver on port " + conf.Server.BindAddress + ":" + strconv.Itoa(conf.Server.BindPort)
	appLogger.Info(logText)
	common.UserLogInfo(logText)

	http.Handle("/", http.FileServer(http.Dir("server/static")))
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/archive", archiveHandler)
	http.HandleFunc("/settings", settingsHandler)
	http.HandleFunc("/view/", viewConfHandler)
	http.HandleFunc("/download/", downloadConfHandler)
	http.HandleFunc("/devicelist", deviceListHandler)
	http.HandleFunc("/devicetypes", deviceTypesHandler)

	err := http.ListenAndServe(conf.Server.BindAddress+":"+strconv.Itoa(conf.Server.BindPort), nil)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	return
}

// Wrapper to render template of name
func renderTemplate(w http.ResponseWriter, name string, d interface{}) {
	if err := templates.ExecuteTemplate(w, name, d); err != nil {
		appLogger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

// Generic function to recover from server errors
func httpRecovery(w http.ResponseWriter) {
	if re := recover(); re != nil {
		appLogger.Errorf("%s", re)
		errorMess := struct{ ErrorMessage string }{"An internal server error has occured."}
		renderTemplate(w, "errorpage", errorMess)
		return
	}
}

// Get a list of all devices in the config.FullConfDir directory
func getDeviceList() deviceList {
	configFileList, _ := ioutil.ReadDir(config.FullConfDir)

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
