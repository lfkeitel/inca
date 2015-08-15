package api

import (
	"net/http"
	"strconv"

	"github.com/dragonrider23/inca/devices"
)

func cpAPI(r *http.Request, urlPieces []string) (interface{}, *apiError) {
	switch urlPieces[0] {
	case "save":
		return "", cpSave(r)
	case "delete":
		return "", cpDelete(r)
	default:
		return "", newError("Endpoint /cp/"+urlPieces[0]+" not found", 1)
	}

	return "Invalid job", nil
}

func cpSave(r *http.Request) *apiError {
	formValues, err := getParams(r,
		[]string{
			"id",
			"name",
			"protocol",
		}, map[string]string{
			"username": "",
			"password": "",
			"enable":   "",
		})
	if err != nil {
		return newError("Make sure all required fields are filled in", 2)
	}

	id, err1 := strconv.Atoi(formValues["id"])
	if err1 != nil {
		id = -1
	}

	cp := devices.ConnProfile{
		Profileid: id,
		Name:      formValues["name"],
		Protocol:  formValues["protocol"],
		Username:  formValues["username"],
		Password:  formValues["password"],
		Enable:    formValues["enable"],
	}

	if id == -1 {
		devices.CreateConnProfile(cp)
	} else {
		devices.EditConnProfile(cp)
	}
	return nil
}

func cpDelete(r *http.Request) *apiError {
	formValues, err := getParams(r,
		[]string{
			"ids",
		}, nil)
	if err != nil {
		return newError("Make sure all required fields are filled in", 2)
	}

	ids, err1 := jsonUnmarshallIntArray(formValues["ids"])
	if err1 != nil {
		return newError(err.Error(), 2)
	}

	err1 = devices.DeleteConnProfiles(ids)
	if err1 != nil {
		return newError(err.Error(), 2)
	}

	return newEmptyError()
}
