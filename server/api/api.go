package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type jsonResponse struct {
	ErrorMessage string
	ErrorCode    int
	Data         interface{}
	Path         string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	urlPieces := strings.Split(r.URL.Path, "/")[2:]
	apiPath := r.URL.Path[len("/api"):]
	var data interface{}
	var err *apiError

	switch urlPieces[0] {
	case "devices":
		data, err = devicesAPI(r, urlPieces[1:])
	// case "inca":
	// 	data, err = incaAPI(r, urlPieces[1:])
	// case "scripts":
	// 	data, err = scriptsAPI(r, urlPieces[1:])
	default:
		err = newError("Module "+urlPieces[0]+" not found", 1)
	}

	if err == nil {
		err = newEmptyError()
	}

	response, _ := prepareResponseJSON(data, err, apiPath)
	w.Write(response)
	return
}

func prepareResponseJSON(d interface{}, e *apiError, p string) ([]byte, error) {
	data := jsonResponse{
		ErrorMessage: e.Error(),
		ErrorCode:    e.Code(),
		Data:         d,
		Path:         p,
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return b, nil
}

func getRequiredParams(r *http.Request, k []string) (map[string]string, error) {
	values := make(map[string]string)

	for _, key := range k {
		v := r.FormValue(key)
		if v == "" {
			return nil, errors.New("Parameter '" + key + "' missing")
		}
		values[key] = v
	}

	return values, nil
}
