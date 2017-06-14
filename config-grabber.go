package main

import (
	"github.com/BurntSushi/toml"

	"github.com/lfkeitel/inca/common"
	"github.com/lfkeitel/inca/grabber"
	"github.com/lfkeitel/inca/server"
	"github.com/lfkeitel/verbose"
)

var appLogger *verbose.Logger

func init() {
	appLogger = verbose.New("app")

	fileLogger, err := verbose.NewFileHandler("logs/app/")
	if err != nil {
		panic("Failed to open logging directory")
	}

	appLogger.AddHandler("file", fileLogger)
}

func main() {
	conf, _ := loadAppConfig()
	grabber.LoadConfig(conf)
	server.StartServer(conf)
}

func loadAppConfig() (common.Config, error) {
	var conf common.Config
	if _, err := toml.DecodeFile("config/configuration.toml", &conf); err != nil {
		appLogger.Fatalf("Couldn't load configuration: %s", err.Error())
		return common.Config{}, err
	}
	return conf, nil
}
