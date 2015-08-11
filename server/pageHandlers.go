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
		// Main management page
		devices, err := devices.GetAllDevices()
		if err != nil {
			httpErrorPage(w, "Failed to load device management page")
			return
		}

		renderTemplate(w, "allDevicesPage", devices)
	} else if deviceID == "new" {
		// Empty device page for new device
		connProfiles, err := devices.GetConnProfiles()
		if err != nil {
			httpErrorPage(w, "Failed to load device management page")
			return
		}

		data := struct {
			Device       devices.Device
			ConnProfiles []devices.ConnProfile
		}{devices.Device{}, connProfiles}
		renderTemplate(w, "singleDevicePage", data)
	} else {
		// Device page to edit properties
		// Convert string to int
		v, err := strconv.Atoi(deviceID)
		if err != nil {
			http.Redirect(w, r, deviceMgt, http.StatusTemporaryRedirect)
			return
		}

		// Get device from database
		device, err := devices.GetDevice(v)
		if err != nil {
			if devices.IsEmptyResultErr(err) {
				http.Redirect(w, r, deviceMgt, http.StatusTemporaryRedirect)
			} else {
				httpErrorPage(w, "Failed to load device management page")
			}
			return
		}

		// Get connection profiles
		connProfiles, err := devices.GetConnProfiles()
		if err != nil {
			httpErrorPage(w, "Failed to load device management page")
			return
		}

		// Build data and render
		data := struct {
			Device       devices.Device
			ConnProfiles []devices.ConnProfile
		}{device, connProfiles}
		renderTemplate(w, "singleDevicePage", data)
	}
}
