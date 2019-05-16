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
		break

	case "runnow":
		response = api.runnow()
		break

	case "singlerun":
		response = api.singlerun(r)
		break

	case "status":
		response = api.status()
		break

	case "devicelist":
		response = api.devicelist()
		break

	case "savedevicelist":
		response = api.savedevicelist(r)
		break

	case "savedevicetypes":
		response = api.savedevicetypes(r)
		break

	case "errorlog":
		response = api.errorlog(r)
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
	splitURL := strings.Split(r.URL.Path, "/")     // [0] = '', [1] = 'view', [2] = filename
	splitName := strings.Split(splitURL[2], "-")   // [0] = name, [1] = datesuffix, [2] = hostname, [3] = dtype
	splitProto := strings.Split(splitName[4], ".") // [0] = method, [1] = ".conf"
	confText, err := ioutil.ReadFile(config.FullConfDir + "/" + splitURL[2])
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
	return
}

// Get the raw configuration file as a download
func downloadConfHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	splitURL := strings.Split(r.URL.Path, "/") // [0] = '', [1] = 'download', [2] = filename
	confText, err := ioutil.ReadFile(config.FullConfDir + "/" + splitURL[2])
	if err != nil {
		panic(err.Error())
	}

	w.Write(confText)
	return
}

// Delete a configuration file
func deleteConfHandler(w http.ResponseWriter, r *http.Request) {
	defer httpRecovery(w)
	splitURL := strings.Split(r.URL.Path, "/") // [0] = '', [1] = 'download', [2] = filename
	os.Remove(filepath.Join(config.FullConfDir, splitURL[2]))
	http.Redirect(w, r, "/archive", http.StatusTemporaryRedirect)
	return
}
