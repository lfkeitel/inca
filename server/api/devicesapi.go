package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dragonrider23/infrastructure-config-archive/devices"
)

func devicesAPI(r *http.Request, urlPieces []string) (interface{}, *apiError) {
	r.ParseForm()

	switch urlPieces[0] {
	case "create":
		return "", create(r)
	}

	return "Invalid job", nil
}

func create(r *http.Request) *apiError {
	formValues, err := getRequiredParams(r,
		[]string{
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

	err = devices.CreateDevice(d)
	if err == nil {
		return newEmptyError()
	}
	return newError(err.Error(), 1)
}

func getRequiredParams(r *http.Request, k []string) (map[string]string, error) {
	var values map[string]string

	for _, key := range k {
		v := r.FormValue(key)
		if v == "" {
			return nil, errors.New("Parameter '" + key + "' missing")
		}
		values[key] = v
	}

	return values, nil
}
