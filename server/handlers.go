package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dragonrider23/infrastructure-config-archive/grabber"
)

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
	splitName := strings.Split(splitUrl[2], "-")   // [0] = name, [1] = datesuffix, [2] = hostname, [3] = dtype
	splitProto := strings.Split(splitName[4], ".") // [0] = method, [1] = ".conf"
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
