package api

import (
	"encoding/json"
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

// Handler dispatches an api request to its appropiate module
func Handler(w http.ResponseWriter, r *http.Request) {
	urlPieces := strings.Split(r.URL.Path, "/")[2:]
	apiPath := r.URL.Path[len("/api"):]
	var data interface{}
	var err *apiError

	switch urlPieces[0] {
	case "devices":
		data, err = devicesAPI(r, urlPieces[1:])
	case "cp":
		data, err = cpAPI(r, urlPieces[1:])
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

func getParams(r *http.Request, req []string, opt map[string]string) (map[string]string, *apiError) {
	r.ParseForm()
	form := r.Form
	values := make(map[string]string, len(opt)+len(req))

	if req != nil {
		for _, key := range req {
			v, ok := form[key]
			if !ok {
				return nil, newError("Parameter '"+key+"' missing", 3)
			}
			values[key] = v[0]
		}
	}

	if opt != nil {
		for key, def := range opt {
			v, ok := form[key]
			if ok {
				v2 := v[0]
				if v2 == "" {
					values[key] = def
				} else {
					values[key] = v2
				}
			} else {
				values[key] = def
			}
		}
	}

	return values, nil
}
