package api

import (
	"net/http"
	"strconv"

	"github.com/dragonrider23/infrastructure-config-archive/devices"
)

func devicesAPI(r *http.Request, urlPieces []string) (interface{}, *apiError) {
	r.ParseForm()

	switch urlPieces[0] {
	case "save":
		return "", save(r)
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
		return newError("Make sure all required fields are filled in", 1)
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
			return newError(err.Error(), 1)
		}
	} else {
		d.Deviceid = id
		err = devices.EditDevice(d)
		if err != nil {
			return newError(err.Error(), 1)
		}
	}
	return newEmptyError()
}
