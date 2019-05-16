package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lfkeitel/inca/src/common"
	"github.com/lfkeitel/inca/src/grabber"
)

type apiRequest struct{}

func (a *apiRequest) running() string {
	return fmt.Sprintf(`{"running": %t}`, grabber.IsRunning())
}

func (a *apiRequest) runnow() string {
	go grabber.PerformConfigGrab()
	return `{"status": "started", "running": true}`
}

func (a *apiRequest) singlerun(r *http.Request) string {
	name := r.FormValue("name")
	hostname := r.FormValue("hostname")
	brand := r.FormValue("brand")
	proto := r.FormValue("proto")
	go grabber.PerformSingleRun(name, hostname, brand, proto)
	return `{"status": "started", "running": true}`
}

func (a *apiRequest) status() string {
	state := grabber.CurrentState()
	return fmt.Sprintf(`{
	"running": %t,
	"totalDevices": %d,
	"finished": %d,
	"stage": "%s"
}`,
		grabber.IsRunning(),
		state.Total,
		state.Finished,
		state.Stage,
	)
}

func (a *apiRequest) devicelist() string {
	deviceList, _ := json.Marshal(getDeviceList())
	return string(deviceList)
}

func (a *apiRequest) savedevicelist(r *http.Request) string {
	listText, _ := url.QueryUnescape(r.FormValue("text"))
	listText = strings.Replace(listText, "-", "_", -1)
	return saveDeviceConfigFile(config.DeviceListFile, listText)
}

func (a *apiRequest) savedevicetypes(r *http.Request) string {
	listText, _ := url.QueryUnescape(r.FormValue("text"))
	return saveDeviceConfigFile(config.DeviceTypeFile, listText)

}

// Save text t to file n after validating the text formatting
func saveDeviceConfigFile(n, t string) string {
	if err := grabber.CheckDeviceList(t); err != nil {
		return fmt.Sprintf(`{"success": false, "error": "%s"}`, err.Error())
	}

	err := ioutil.WriteFile(n, []byte(t), 0664)
	if err != nil {
		return fmt.Sprintf(`{"success": false, "error": "%s"}`, err.Error())
	}
	return `{"success": true}`
}

type errorLogLine struct {
	Etype   string `json:"etype"`
	Time    string `json:"time"`
	Message string `json:"message"`
}

func (a *apiRequest) errorlog(r *http.Request) string {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	log, err := ioutil.ReadFile("logs/endUser.log")
	if err != nil {
		return "{}"
	}

	logLines := strings.Split(string(log), "\n")

	// Remove last element if blank line
	if logLines[len(logLines)-1] == "" {
		logLines = logLines[:len(logLines)-1]
	}

	logLines = common.ReverseSlice(logLines)

	// If the slice is longer than the requested events, shorten it
	if len(logLines) > limit {
		logLines = append([]string(nil), logLines[:limit]...)
	}

	// Parse the log lines into their elemental parts
	parsedLines := []errorLogLine{}
	for _, l := range logLines {
		line := strings.Split(l, ": ")

		sLine := errorLogLine{
			Etype:   line[1],
			Time:    line[0],
			Message: line[3],
		}
		parsedLines = append(parsedLines, sLine)
	}

	logJSON, err := json.Marshal(parsedLines)
	if err != nil {
		return "{}"
	}
	return string(logJSON)
}
