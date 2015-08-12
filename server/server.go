package server

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/dragonrider23/go-logger"
	"github.com/dragonrider23/inca/common"
	"github.com/dragonrider23/inca/server/api"
)

const (
	rootPath      = "/"
	staticContent = rootPath + "static/"
	deviceMgt     = rootPath + "devices/"
	apiPath       = rootPath + "api/"
)

var templates *template.Template
var appLogger *logger.Logger
var config common.Config

// Initialize HTTP server with app configuration and templates
func initServer(configuration common.Config) {
	config = configuration
	templates = template.Must(template.ParseGlob("server/templates/*.tmpl"))
	appLogger = logger.New("httpServer").Path("logs/server/")
}

// Start Start front-end HTTP server
func Start(conf common.Config) {
	initServer(conf)

	logText := "Starting webserver on port " + conf.Server.BindAddress + ":" + strconv.Itoa(conf.Server.BindPort)
	appLogger.Verbose(3)
	appLogger.Info(logText)
	common.UserLogInfo(logText)

	http.Handle(staticContent, http.FileServer(http.Dir("server")))
	http.HandleFunc(rootPath, indexHandler)
	http.HandleFunc(deviceMgt, deviceMgtHandler)
	http.HandleFunc(apiPath, api.Handler)

	err := http.ListenAndServe(conf.Server.BindAddress+":"+strconv.Itoa(conf.Server.BindPort), nil)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	return
}

// Wrapper to render template of name
func renderTemplate(w http.ResponseWriter, name string, d interface{}) {
	err := templates.ExecuteTemplate(w, name, d)
	if isErr := logger.CheckError(err, appLogger); isErr {
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

func httpErrorPage(w http.ResponseWriter, msg string) {
	appLogger.Error("%s", msg)
	errorMess := struct{ ErrorMessage string }{msg}
	renderTemplate(w, "errorpage", errorMess)
	return
}
