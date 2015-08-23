package api

import (
	"net/http"
	"strconv"

	"github.com/dragonrider23/inca/devices"
	"github.com/dragonrider23/inca/poller"
)

func devicesAPI(r *http.Request, urlPieces []string) (interface{}, *apiError) {
	r.ParseForm()

	switch urlPieces[0] {
	case "save":
		return "", deviceSave(r)
	case "delete":
		return "", deviceDelete(r)
	case "update":
		return update(r)
	case "status":
		return deviceStatus(r)
	default:
		return "", newError("Endpoint /devices/"+urlPieces[0]+" not found", 1)
	}

	return "Invalid job", nil
}

func deviceSave(r *http.Request) *apiError {
	formValues, err := getParams(r,
		[]string{
			"name",
			"hostname",
			"connprofile",
			"disabled",
		}, map[string]string{
			"deviceid":     "-1",
			"manufacturer": "",
			"model":        "",
		})
	if err != nil {
		return newError("Make sure all required fields are filled in", 2)
	}

	cp, err1 := strconv.Atoi(formValues["connprofile"])
	if err1 != nil {
		cp = 0
	}

	id, err1 := strconv.Atoi(formValues["deviceid"])
	if err1 != nil {
		id = -1
	}

	disabled, err1 := strconv.ParseBool(formValues["disabled"])
	if err1 != nil {
		disabled = false
	}

	d := devices.Device{
		Name:         formValues["name"],
		Hostname:     formValues["hostname"],
		ConnProfile:  cp,
		Manufacturer: formValues["manufacturer"],
		Model:        formValues["model"],
		Disabled:     disabled,
	}

	if id == -1 {
		err1 = devices.CreateDevice(d)
		if err1 != nil {
			return newError(err.Error(), 2)
		}
	} else {
		d.Deviceid = id
		err1 = devices.EditDevice(d)
		if err1 != nil {
			return newError(err.Error(), 2)
		}
	}
	return newEmptyError()
}

func deviceDelete(r *http.Request) *apiError {
	formValues, err := getParams(r,
		[]string{
			"deviceids",
		}, nil)
	if err != nil {
		return newError("Make sure all required fields are filled in", 2)
	}

	ids, err1 := jsonUnmarshallIntArray(formValues["deviceids"])
	if err1 != nil {
		return newError(err.Error(), 2)
	}

	err1 = devices.DeleteDevices(ids)
	if err1 != nil {
		return newError(err.Error(), 2)
	}

	return newEmptyError()
}

func update(r *http.Request) (string, *apiError) {
	formValues, er := getParams(r, nil, map[string]string{
		"deviceids": "-1",
	})
	if er != nil {
		return "", newError("Make sure all required fields are filled in", 2)
	}

	var d []devices.Device
	var err error

	if formValues["deviceids"] == "-1" {
		d, err = devices.GetDevicesForConfigGrab(nil)
		if err != nil {
			return "", newError("Failed to start poll job x1", 2)
		}
	} else {
		ids, err := jsonUnmarshallIntArray(formValues["deviceids"])
		if err != nil {
			return "", newError(err.Error(), 2)
		}

		d, err = devices.GetDevicesForConfigGrab(ids)
		if err != nil {
			return "", newError("Failed to start poll job x2", 2)
		}
	}

	cp, err := devices.GetConnProfiles()
	if err != nil {
		return "", newError("Failed to start poll job x3", 2)
	}

	_, _, err = poller.Process(poller.Job{
		Cmd: "poll",
		Data: map[string]interface{}{
			"devices":      d,
			"connProfiles": cp,
		},
	})

	if err != nil {
		return "", newError(err.Error(), 1)
	}

	return "Update started", newEmptyError()
}

func deviceStatus(r *http.Request) (devices.DeviceStatus, *apiError) {
	d, err := devices.GetDeviceStats()
	if err != nil {
		return d, newError(err.Error(), 2)
	}
	return d, newEmptyError()
}
