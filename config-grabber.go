package main

import (
	"github.com/BurntSushi/toml"

	logger "github.com/dragonrider23/go-logger"
	"github.com/dragonrider23/infrastructure-config-archive/comm"
	"github.com/dragonrider23/infrastructure-config-archive/grabber"
	"github.com/dragonrider23/infrastructure-config-archive/server"
)

var appLogger = logger.New("app").Verbose(3).Path("logs/app/")

func loadAppConfig() (comm.Config, error) {
	var conf comm.Config
	if _, err := toml.DecodeFile("config/configuration.toml", &conf); err != nil {
		appLogger.Fatal("Couldn't load configuration: %s", err.Error())
		return comm.Config{}, err
	}
	return conf, nil
}

func main() {
	conf, _ := loadAppConfig()
	grabber.LoadConfig(conf)
	server.StartServer(conf)
	return
}
