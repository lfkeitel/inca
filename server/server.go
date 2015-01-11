package server

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/dragonrider23/config-grabber/interfaces"
	"github.com/dragonrider23/config-grabber/grabber"
	logs "github.com/dragonrider23/go-logger"
)

type deviceConfigFile struct {
	Path string
	Name string
	Address string
	Proto string
	ConfText []string
}

type deviceList struct {
	Devices []deviceConfigFile
}

var templates *template.Template
var appLogger *logs.Logger
var config interfaces.Config

// Initialize HTTP server with app configuration and templates
func initServer(configuration interfaces.Config) {
	config = configuration
	templates = template.Must(template.ParseGlob(config.Server.BaseDir+"/templates/*.tmpl"))
	appLogger = logs.New("httpServer")
}

// Wrapper to render template of name
func renderTemplate(w http.ResponseWriter, name string, d interface{}) {
	err := templates.ExecuteTemplate(w, name, d)
	if isErr := logs.CheckError(err, appLogger); isErr {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

// Generic function to recover from server errors
func httpRecovery(w http.ResponseWriter) {
	if re := recover(); re != nil {
		appLogger.Error("%s", re)
		errorMess := struct{ ErrorMessage string }{"An internal server error has occured."}
		renderTemplate(w, "errorpage", errorMess)
		return
	}
}

// Handler for general API functions such as status, and initiating a config run
func apiHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)

	splitUrl := strings.Split(r.URL.Path, "/") // [0] = '', [1] = 'api', [2:] = Custom path
	response := ""

	switch splitUrl[2] {
		case "running":
			response = "{\"status\": "+strconv.FormatBool(grabber.IsRunning())+" }"
			break
		case "runnow":
			go grabber.PerformConfigGrab()
			response = "{\"status\": \"started\" }"
			break
		case "status":
			total, finished := grabber.Remaining()
			response = "{\"status\": "+strconv.FormatBool(grabber.IsRunning())+", \"totalDevices\": "+strconv.Itoa(total)+", \"finished\": "+strconv.Itoa(finished)+"}"
			break
		case "devicelist":
			deviceList, _ := json.Marshal(getDeviceList())
			response = string(deviceList)
			break
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
	return
}

// Generate page with a list of device configurations
func archiveHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	renderTemplate(w, "archivePage", getDeviceList())
	return
}

// Generate page to display configuration of given file (in URL)
func viewConfHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	splitUrl := strings.Split(r.URL.Path, "/") // [0] = '', [1] = 'view', [2] = filename
	splitName := strings.Split(splitUrl[2], "-") // [0] = name, [1] = datesuffix, [2] = hostname
	splitProto := strings.Split(splitName[3], ".") // [0] = protocol, [1] = ".conf"
	confText, err := ioutil.ReadFile(config.FullConfDir+"/"+splitUrl[2])
	if err != nil {
		panic(err.Error())
	}

	device := deviceConfigFile{
		Path: splitUrl[2],
		Name: splitName[0],
		Address: splitName[2],
		Proto: splitProto[0],
		ConfText: strings.Split(string(confText), "\n"),
	}

	renderTemplate(w, "viewConfPage", device)
	return
}

// Get the raw configuration file as a download
func downloadConfHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	splitUrl := strings.Split(r.URL.Path, "/") // [0] = '', [1] = 'download', [2] = filename
	confText, err := ioutil.ReadFile(config.FullConfDir+"/"+splitUrl[2])
	if err != nil {
		panic(err.Error())
	}

	w.Write(confText)
	return
}

// Get a list of all devices in the configuration directory
func getDeviceList() deviceList {
	configFileList, _ := ioutil.ReadDir(config.FullConfDir);

	deviceConfigs := deviceList{}

	for _, file := range configFileList {
		splitName := strings.Split(file.Name(), "-") // [0] = name, [1] = datesuffix, [2] = hostname
		splitProto := strings.Split(splitName[3], ".") // [0] = protocol, [1] = ".conf"

		device := deviceConfigFile{
			Path: file.Name(),
			Name: splitName[0],
			Address: splitName[2],
			Proto: splitProto[0],
		}
		deviceConfigs.Devices = append(deviceConfigs.Devices, device)
	}

	return deviceConfigs
}

// Start front-end HTTP server
func StartServer(conf interfaces.Config) {
	initServer(conf)

	appLogger.Verbose(3)
	appLogger.Info("Starting webserver on port "+conf.Server.BindAddress+":"+strconv.Itoa(conf.Server.BindPort))

	http.Handle("/", http.FileServer(http.Dir(conf.Server.BaseDir+"/static")))
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/archive", archiveHandler)
	http.HandleFunc("/view/", viewConfHandler)
	http.HandleFunc("/download/", downloadConfHandler)

	appLogger.Info("Server ready")
	err := http.ListenAndServe(conf.Server.BindAddress+":"+strconv.Itoa(conf.Server.BindPort), nil)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	return
}
