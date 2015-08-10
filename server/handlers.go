package server

import (
	"net/http"
	"strconv"

	"github.com/dragonrider23/infrastructure-config-archive/devices"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == rootPath {
		renderTemplate(w, "index", nil)
	} else {
		http.NotFound(w, r)
	}
}

func deviceMgtHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Path[len(deviceMgt):]

	if deviceID == "" {
		devices, err := devices.GetAllDevices()
		if err == nil {
			renderTemplate(w, "allDevicesPage", devices)
		} else {
			httpErrorPage(w, "Failed to load device management page")
		}
	} else if deviceID == "new" {
		connProfiles, err := devices.GetConnProfiles()
		if err != nil {
			httpErrorPage(w, "Failed to load device management page")
		} else {
			data := struct {
				Device       devices.Device
				ConnProfiles []devices.ConnProfile
			}{devices.Device{}, connProfiles}
			renderTemplate(w, "singleDevicePage", data)
		}
	} else {
		v, err := strconv.Atoi(deviceID)
		if err == nil {
			device, err := devices.GetDevice(v)
			if err != nil {
				httpErrorPage(w, "Failed to load device management page")
			} else {
				connProfiles, err := devices.GetConnProfiles()
				if err != nil {
					httpErrorPage(w, "Failed to load device management page")
				} else {
					data := struct {
						Device       devices.Device
						ConnProfiles []devices.ConnProfile
					}{device, connProfiles}
					renderTemplate(w, "singleDevicePage", data)
				}
			}
		} else {
			http.Redirect(w, r, deviceMgt, http.StatusTemporaryRedirect)
		}
	}
}
