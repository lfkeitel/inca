package main

import (
	"github.com/BurntSushi/toml"

	"github.com/dragonrider23/config-grabber/server"
	"github.com/dragonrider23/config-grabber/interfaces"
    "github.com/dragonrider23/config-grabber/grabber"
	logger "github.com/dragonrider23/go-logger"
)

var appLogger = logger.New("app").Verbose(3)

func loadAppConfig() (interfaces.Config, error) {
	var conf interfaces.Config
	if _, err := toml.DecodeFile("configuration.toml", &conf); err != nil {
		appLogger.Fatal("Couldn't load configuration: %s", err.Error())
		return interfaces.Config{}, err
	}
	return conf, nil
}

func main() {
	conf, _ := loadAppConfig()
	grabber.LoadConfig(conf)
	server.StartServer(conf)
	return
}
