package server

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dragonrider23/inca/internal/common"
	"github.com/dragonrider23/inca/logger"
	"github.com/dragonrider23/inca/server/api"
)

const (
	rootPath      = "/"
	staticContent = rootPath + "static/"
	deviceMgt     = rootPath + "devices/"
	apiPath       = rootPath + "api/"
	adminPage     = rootPath + "admin/"
)

var templates *template.Template
var appLogger *logger.Logger

// Initialize HTTP server with app configuration and templates
func initServer() {
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"dict": func(v ...interface{}) (map[string]interface{}, error) {
			if len(v)%2 != 0 {
				return nil, errors.New("Invalid dict call")
			}
			dict := make(map[string]interface{}, len(v)/2)
			for i := 0; i < len(v); i += 2 {
				key, ok := v[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = v[i+1]
			}
			return dict, nil
		},

		"list": func(v ...interface{}) ([]interface{}, error) {
			list := make([]interface{}, len(v))
			for i := 0; i < len(v); i++ {
				list[i] = v[i]
			}
			return list, nil
		},
	}).ParseGlob("server/templates/*/*.tmpl"))

	appLogger = logger.New("webserver")
}

// Start Start front-end HTTP server
func Start() {
	initServer()
	conf := common.Config

	logText := "Starting webserver on port " + conf.Server.BindAddress + ":" + strconv.Itoa(conf.Server.BindPort)
	appLogger.Info(logText)

	http.Handle(staticContent, http.FileServer(http.Dir("server")))
	http.HandleFunc(rootPath, indexHandler)
	http.HandleFunc(deviceMgt, deviceMgtHandler)
	http.HandleFunc(apiPath, api.Handler)
	http.HandleFunc(adminPage, adminPageHandler)

	err := http.ListenAndServe(conf.Server.BindAddress+":"+strconv.Itoa(conf.Server.BindPort), nil)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	return
}

// Wrapper to render template of name
func renderTemplate(w http.ResponseWriter, name string, d interface{}) {
	err := templates.ExecuteTemplate(w, name, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

// Generic function to recover from server errors
func httpRecovery(w http.ResponseWriter) {
	if re := recover(); re != nil {
		appLogger.Error("%s", re)
		errorMess := struct{ ErrorMessage string }{"An internal server error has occured."}
		renderTemplate(w, "errorpage", errorMess)
		return
	}
}

func httpErrorPage(w http.ResponseWriter, msg string, log bool) {
	if log {
		appLogger.Error("%s", msg)
	}
	errorMess := struct{ ErrorMessage string }{msg}
	renderTemplate(w, "errorpage", errorMess)
	return
}
