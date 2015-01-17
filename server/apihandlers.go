package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dragonrider23/infrastructure-config-archive/grabber"
)

type apiRequest struct{}

func (a *apiRequest) running() string {
	return "{\"running\": " + strconv.FormatBool(grabber.IsRunning()) + " }"
}

func (a *apiRequest) runnow() string {
	go grabber.PerformConfigGrab()
	return "{\"status\": \"started\", \"running\": true}"
}

func (a *apiRequest) singlerun(r *http.Request) string {
	name := r.FormValue("name")
	hostname := r.FormValue("hostname")
	brand := r.FormValue("brand")
	proto := r.FormValue("proto")
	go grabber.PerformSingleRun(name, hostname, brand, proto)
	return "{\"status\": \"started\", \"running\": true}"
}

func (a *apiRequest) status() string {
	total, finished := grabber.Remaining()
	return "{\"status\": " + strconv.FormatBool(grabber.IsRunning()) + ", \"running\": " + strconv.FormatBool(grabber.IsRunning()) + ", \"totalDevices\": " + strconv.Itoa(total) + ", \"finished\": " + strconv.Itoa(finished) + "}"
}

func (a *apiRequest) devicelist() string {
	deviceList, _ := json.Marshal(getDeviceList())
	return string(deviceList)
}

func (a *apiRequest) savedevicelist(r *http.Request) string {
	listText, _ := url.QueryUnescape(r.FormValue("text"))
	err := ioutil.WriteFile(config.DeviceListFile, []byte(listText), 0664)
	if err != nil {
		return "{\"success\": false, \"error\": \"" + err.Error() + "\"}"
	} else {
		return "{\"success\": true}"
	}
}

func (a *apiRequest) savedevicetypes(r *http.Request) string {
	listText, _ := url.QueryUnescape(r.FormValue("text"))
	err := ioutil.WriteFile(config.DeviceTypeFile, []byte(listText), 0664)
	if err != nil {
		return "{\"success\": false, \"error\": \"" + err.Error() + "\"}"
	} else {
		return "{\"success\": true}"
	}
}
