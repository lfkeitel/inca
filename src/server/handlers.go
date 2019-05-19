package server

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler for general API functions such as status, and initiating a config run
func apiHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)

	splitURL := strings.Split(r.URL.Path, "/") // [0] = '', [1] = 'api', [2:] = Custom path
	var response string
	api := apiRequest{}

	switch splitURL[2] {
	case "running":
		response = api.running()
	case "runnow":
		response = api.runnow()
	case "singlerun":
		response = api.singlerun(r)
	case "status":
		response = api.status()
	case "devicelist":
		response = api.devicelist()
	case "savedevicelist":
		response = api.savedevicelist(r)
	case "savedevicetypes":
		response = api.savedevicetypes(r)
	case "errorlog":
		response = api.errorlog(r)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}

// Generate page with a list of device configurations
func archiveHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	renderTemplate(w, "archivePage", getDeviceList())
}

// Generate page with the device definitions
func deviceListHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	confText, err := ioutil.ReadFile(config.Paths.DeviceList)
	if err != nil {
		panic(err.Error())
	}

	data := struct {
		ConfText string
		Path     string
	}{string(confText), config.Paths.DeviceList}
	renderTemplate(w, "deviceListPage", data)
}

// Generate page with the device type definitions
func deviceTypesHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	confText, err := ioutil.ReadFile(config.Paths.DeviceTypes)
	if err != nil {
		panic(err.Error())
	}

	data := struct {
		ConfText string
		Path     string
	}{string(confText), config.Paths.DeviceTypes}
	renderTemplate(w, "deviceTypePage", data)
}

// Generate page to display configuration of given file (in URL)
func viewConfHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	splitURL := strings.Split(r.URL.Path, "/")     // [0] = '', [1] = 'view', [2] = filename
	splitName := strings.Split(splitURL[2], "-")   // [0] = name, [1] = datesuffix, [2] = hostname, [3] = dtype
	splitProto := strings.Split(splitName[4], ".") // [0] = method, [1] = ".conf"
	confText, err := ioutil.ReadFile(filepath.Join(config.Paths.ConfDir, splitURL[2]))
	if err != nil {
		panic(err.Error())
	}

	device := deviceConfigFile{
		Path:         splitURL[2],
		Name:         splitName[0],
		Address:      splitName[2],
		Proto:        splitProto[0],
		ConfText:     strings.Split(string(confText), "\n"),
		Manufacturer: splitName[3],
	}

	renderTemplate(w, "viewConfPage", device)
}

// Get the raw configuration file as a download
func downloadConfHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	splitURL := strings.Split(r.URL.Path, "/") // [0] = '', [1] = 'download', [2] = filename
	confText, err := ioutil.ReadFile(filepath.Join(config.Paths.ConfDir, splitURL[2]))
	if err != nil {
		panic(err.Error())
	}

	w.Write(confText)
}

// Delete a configuration file
func deleteConfHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	splitURL := strings.Split(r.URL.Path, "/") // [0] = '', [1] = 'download', [2] = filename
	os.Remove(filepath.Join(config.Paths.ConfDir, splitURL[2]))
	http.Redirect(w, r, "/archive", http.StatusTemporaryRedirect)
}
