package server

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// Handler for general API functions such as status, and initiating a config run
func apiHandler(w http.ResponseWriter, r *http.Request) {
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
	case "download":
		downloadConfApiHandler(w, r)
		return
	case "getdevicelistfile":
		configFileApiHandler(w, r, config.Paths.DeviceList)
		return
	case "getdevicetypesfile":
		configFileApiHandler(w, r, config.Paths.DeviceTypes)
		return
	case "delete":
		response = api.deleteconf(r)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}

// Generate page with the device type definitions
func configFileApiHandler(w http.ResponseWriter, r *http.Request, file string) {
	confText, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-Type", "plain/text")
	w.Write(confText)
}

func downloadConfApiHandler(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.Path, "/")
	path := splitURL[3]

	w.Header().Set("Content-Type", "plain/text")

	confText, err := ioutil.ReadFile(filepath.Join(config.Paths.ConfDir, path))
	if err != nil {
		w.Write([]byte("Configuration file does not exist"))
		return
	}

	w.Write(confText)
}
