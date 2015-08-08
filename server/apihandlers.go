package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	//"net/url"
	"strconv"
	"strings"

	"github.com/dragonrider23/utils/slices"

	"github.com/dragonrider23/infrastructure-config-archive/configs"
)

type apiRequest struct{}

func (a *apiRequest) running() string {
	return "{\"running\": " + strconv.FormatBool(configs.IsRunning()) + " }"
}

func (a *apiRequest) runnow() string {
	go configs.PerformConfigGrab()
	return "{\"status\": \"started\", \"running\": true}"
}

// func (a *apiRequest) singlerun(r *http.Request) string {
// 	// name := r.FormValue("name")
// 	// hostname := r.FormValue("hostname")
// 	// brand := r.FormValue("brand")
// 	// proto := r.FormValue("proto")
// 	//go configs.PerformSingleRun(name, hostname, brand, proto)
// 	return "{\"status\": \"started\", \"running\": true}"
// }

func (a *apiRequest) status() string {
	total, finished := configs.Remaining()
	return "{\"status\": " + strconv.FormatBool(configs.IsRunning()) + ", \"running\": " + strconv.FormatBool(configs.IsRunning()) + ", \"totalDevices\": " + strconv.Itoa(total) + ", \"finished\": " + strconv.Itoa(finished) + "}"
}

// func (a *apiRequest) devicelist() string {
// 	deviceList, _ := json.Marshal(getDeviceList())
// 	return string(deviceList)
// }
//
// func (a *apiRequest) savedevicelist(r *http.Request) string {
// 	listText, _ := url.QueryUnescape(r.FormValue("text"))
// 	return saveDeviceConfigFile(config.DeviceListFile, listText)
// }
//
// func (a *apiRequest) savedevicetypes(r *http.Request) string {
// 	listText, _ := url.QueryUnescape(r.FormValue("text"))
// 	return saveDeviceConfigFile(config.DeviceTypeFile, listText)
//
// }
//
// // Save text t to file n after validating the text formatting
// func saveDeviceConfigFile(n, t string) string {
// 	// if err := configs.CheckDeviceList(t); err != nil {
// 	// 	return "{\"success\": false, \"error\": \"" + err.Error() + "\"}"
// 	// }
//
// 	t = strings.Replace(t, "-", "_", -1)
// 	err := ioutil.WriteFile(n, []byte(t), 0664)
// 	if err != nil {
// 		return "{\"success\": false, \"error\": \"" + err.Error() + "\"}"
// 	}
// 	return "{\"success\": true}"
// }

type errorLogLine struct {
	Etype   string
	Time    string
	Message string
}

func (a *apiRequest) errorlog(r *http.Request) string {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	log, err := ioutil.ReadFile("logs/endUser/enduserlog-log.log")
	if err != nil {
		return "{}"
	}

	logLines := strings.Split(string(log), "\n")

	// Remove last element if blank line
	if logLines[len(logLines)-1] == "" {
		logLines = logLines[:len(logLines)-1]
	}

	logLines = slices.ReverseSlice(logLines)

	// If the slice is longer than the requested events, shorten it
	if len(logLines) > limit {
		logLines = append([]string(nil), logLines[:limit]...)
	}

	// Parse the log lines into their elemental parts
	parsedLines := []errorLogLine{}
	for _, l := range logLines {
		line := strings.Split(l, ":-:")

		sLine := errorLogLine{
			Etype:   line[0],
			Time:    line[1],
			Message: line[2],
		}
		parsedLines = append(parsedLines, sLine)
	}

	logJSON, err := json.Marshal(parsedLines)
	if err != nil {
		return "{}"
	}
	return string(logJSON)
}
