package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	logs "github.com/dragonrider23/go-logger"
	"github.com/dragonrider23/infrastructure-config-archive/grabber"
	"github.com/dragonrider23/infrastructure-config-archive/interfaces"
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
var appLogger *logs.Logger
var config interfaces.Config

// Initialize HTTP server with app configuration and templates
func initServer(configuration interfaces.Config) {
	config = configuration
	templates = template.Must(template.ParseGlob("server/templates/*.tmpl"))
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
		response = "{\"running\": " + strconv.FormatBool(grabber.IsRunning()) + " }"
		break

	case "runnow":
		go grabber.PerformConfigGrab()
		response = "{\"status\": \"started\", \"running\": true}"
		break

	case "singlerun":
		name := r.FormValue("name")
		hostname := r.FormValue("hostname")
		brand := r.FormValue("brand")
		proto := r.FormValue("proto")
		go grabber.PerformSingleRun(name, hostname, brand, proto)
		response = "{\"status\": \"started\", \"running\": true}"
		break

	case "status":
		total, finished := grabber.Remaining()
		response = "{\"status\": " + strconv.FormatBool(grabber.IsRunning()) + ", \"running\": " + strconv.FormatBool(grabber.IsRunning()) + ", \"totalDevices\": " + strconv.Itoa(total) + ", \"finished\": " + strconv.Itoa(finished) + "}"
		break

	case "devicelist":
		deviceList, _ := json.Marshal(getDeviceList())
		response = string(deviceList)
		break

	case "savedevicelist":
		listText, _ := url.QueryUnescape(r.FormValue("text"))
		err := ioutil.WriteFile(config.DeviceListFile, []byte(listText), 0664)
		if err != nil {
			response = "{\"success\": false, \"error\": \"" + err.Error() + "\"}"
		} else {
			response = "{\"success\": true}"
		}
		break

	case "savedevicetypes":
		listText, _ := url.QueryUnescape(r.FormValue("text"))
		err := ioutil.WriteFile(config.DeviceTypeFile, []byte(listText), 0664)
		if err != nil {
			response = "{\"success\": false, \"error\": \"" + err.Error() + "\"}"
		} else {
			response = "{\"success\": true}"
		}
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

// Generate page with the application configuration
func settingsHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	confText, err := ioutil.ReadFile("config/configuration.toml")
	if err != nil {
		panic(err.Error())
	}

	data := struct {
		ConfText []string
	}{strings.Split(string(confText), "\n")}
	renderTemplate(w, "settingsPage", data)
	return
}

// Generate page with the device definitions
func deviceListHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	confText, err := ioutil.ReadFile(config.DeviceListFile)
	if err != nil {
		panic(err.Error())
	}

	data := struct {
		ConfText string
		Path     string
	}{string(confText), config.DeviceListFile}
	renderTemplate(w, "deviceListPage", data)
	return
}

// Generate page with the device type definitions
func deviceTypesHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	confText, err := ioutil.ReadFile(config.DeviceTypeFile)
	if err != nil {
		panic(err.Error())
	}

	data := struct {
		ConfText string
		Path     string
	}{string(confText), config.DeviceTypeFile}
	renderTemplate(w, "deviceTypePage", data)
	return
}

// Generate page to display configuration of given file (in URL)
func viewConfHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	splitUrl := strings.Split(r.URL.Path, "/")     // [0] = '', [1] = 'view', [2] = filename
	splitName := strings.Split(splitUrl[2], "-")   // [0] = name, [1] = datesuffix, [2] = hostname, [3] = manufacturer
	splitProto := strings.Split(splitName[4], ".") // [0] = protocol, [1] = ".conf"
	confText, err := ioutil.ReadFile(config.FullConfDir + "/" + splitUrl[2])
	if err != nil {
		panic(err.Error())
	}

	device := deviceConfigFile{
		Path:         splitUrl[2],
		Name:         splitName[0],
		Address:      splitName[2],
		Proto:        splitProto[0],
		ConfText:     strings.Split(string(confText), "\n"),
		Manufacturer: splitName[3],
	}

	renderTemplate(w, "viewConfPage", device)
	return
}

// Get the raw configuration file as a download
func downloadConfHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	splitUrl := strings.Split(r.URL.Path, "/") // [0] = '', [1] = 'download', [2] = filename
	confText, err := ioutil.ReadFile(config.FullConfDir + "/" + splitUrl[2])
	if err != nil {
		panic(err.Error())
	}

	w.Write(confText)
	return
}

// Get a list of all devices in the configuration directory
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

// Start front-end HTTP server
func StartServer(conf interfaces.Config) {
	initServer(conf)

	appLogger.Verbose(3)
	appLogger.Info("Starting webserver on port " + conf.Server.BindAddress + ":" + strconv.Itoa(conf.Server.BindPort))

	http.Handle("/", http.FileServer(http.Dir("server/static")))
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/archive", archiveHandler)
	http.HandleFunc("/settings", settingsHandler)
	http.HandleFunc("/view/", viewConfHandler)
	http.HandleFunc("/download/", downloadConfHandler)
	http.HandleFunc("/devicelist", deviceListHandler)
	http.HandleFunc("/devicetypes", deviceTypesHandler)

	appLogger.Info("Server ready")
	err := http.ListenAndServe(conf.Server.BindAddress+":"+strconv.Itoa(conf.Server.BindPort), nil)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	return
}
