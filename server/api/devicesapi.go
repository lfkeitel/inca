package api

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"strconv"

	"github.com/dragonrider23/inca/devices"
)

func devicesAPI(r *http.Request, urlPieces []string) (interface{}, *apiError) {
	r.ParseForm()

	switch urlPieces[0] {
	case "save":
		return "", save(r)
	case "delete":
		return "", delete(r)
	// case "update":
	// 	return "", update(r)
	case "status":
		return status(r)
	default:
		return "", newError("Endpoint /devices/"+urlPieces[0]+" not found", 1)
	}

	return "Invalid job", nil
}

func save(r *http.Request) *apiError {
	formValues, err := getRequiredParams(r,
		[]string{
			"deviceid",
			"name",
			"hostname",
			"connprofile",
			"manufacturer",
			"model",
			"disabled",
		})
	if err != nil {
		return newError("Make sure all required fields are filled in", 2)
	}

	cp, err := strconv.Atoi(formValues["connprofile"])
	if err != nil {
		cp = 0
	}

	id, err := strconv.Atoi(formValues["deviceid"])
	if err != nil {
		id = -1
	}

	disabled, err := strconv.ParseBool(formValues["disabled"])
	if err != nil {
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
		err = devices.CreateDevice(d)
		if err != nil {
			return newError(err.Error(), 2)
		}
	} else {
		d.Deviceid = id
		err = devices.EditDevice(d)
		if err != nil {
			return newError(err.Error(), 2)
		}
	}
	return newEmptyError()
}

func delete(r *http.Request) *apiError {
	formValues, err := getRequiredParams(r,
		[]string{
			"deviceids",
		})
	if err != nil {
		return newError("Make sure all required fields are filled in", 2)
	}

	ids, err := jsonUnmarshallDeviceIDs(formValues["deviceids"])
	if err != nil {
		return newError(err.Error(), 2)
	}

	err = devices.DeleteDevices(ids)
	if err != nil {
		return newError(err.Error(), 2)
	}

	return newEmptyError()
}

func update(r *http.Request) *apiError {
	// formValues, err := getRequiredParams(r,
	// 	[]string{
	// 		"deviceids",
	// 	})
	// if err != nil {
	// 	return newError("Make sure all required fields are filled in", 2)
	// }
	//
	// ids, err := jsonUnmarshallDeviceIDs(formValues["deviceids"])
	// if err != nil {
	// 	return newError(err.Error(), 2)
	// }

	// err = devices.DeleteDevices(ids)
	// if err != nil {
	// 	return newError(err.Error(), 2)
	// }
	return newEmptyError()
}

func jsonUnmarshallDeviceIDs(s string) ([]int, error) {
	var ids []int
	if err := json.Unmarshal([]byte(s), &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func status(r *http.Request) (devices.DeviceStatus, *apiError) {
	d, err := devices.GetDeviceStats()
	if err != nil {
		return d, newError(err.Error(), 2)
	}
	return d, newEmptyError()
}
