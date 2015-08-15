package server

import (
	"net/http"
	"strconv"

	"github.com/dragonrider23/inca/devices"
)

func deviceMgtHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Path[len(deviceMgt):]

	if deviceID == "" {
		// Main management page
		query := r.FormValue("query")
		var err error
		var hosts []devices.Device

		if query == "" {
			hosts, err = devices.GetAllDevices()
		} else {
			hosts, err = devices.Search(query)
		}

		if err != nil {
			appLogger.Error(err.Error())
			httpErrorPage(w, "Failed to load device management page", false)
			return
		}

		data := struct {
			Hosts []devices.Device
			Query string
		}{hosts, query}

		renderTemplate(w, "allDevicesPage", data)
	} else if deviceID == "new" {
		// Empty device page for new device
		connProfiles, err := devices.GetConnProfiles()
		if err != nil {
			appLogger.Error(err.Error())
			httpErrorPage(w, "Failed to load device management page", false)
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
				appLogger.Error(err.Error())
				httpErrorPage(w, "Failed to load device management page", false)
			}
			return
		}

		// Get connection profiles
		connProfiles, err := devices.GetConnProfiles()
		if err != nil {
			appLogger.Error(err.Error())
			httpErrorPage(w, "Failed to load device management page", false)
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
