package main

import (
	"github.com/BurntSushi/toml"

	"github.com/lfkeitel/go-logger"
	"github.com/lfkeitel/inca/comm"
	"github.com/lfkeitel/inca/grabber"
	"github.com/lfkeitel/inca/server"
)

var appLogger *logger.Logger

func loadAppConfig() (comm.Config, error) {
	var conf comm.Config
	if _, err := toml.DecodeFile("config/configuration.toml", &conf); err != nil {
		appLogger.Fatal("Couldn't load configuration: %s", err.Error())
		return comm.Config{}, err
	}
	return conf, nil
}

func init() {
	appLogger = logger.New("app").Verbose(3).Path("logs/app/")
}

func main() {
	conf, _ := loadAppConfig()
	grabber.LoadConfig(conf)
	server.StartServer(conf)
	return
}
