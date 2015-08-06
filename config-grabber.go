package main

import (
	"github.com/BurntSushi/toml"

	"github.com/dragonrider23/go-logger"
	"github.com/dragonrider23/infrastructure-config-archive/common"
	"github.com/dragonrider23/infrastructure-config-archive/grabber"
	"github.com/dragonrider23/infrastructure-config-archive/server"
)

var appLogger *logger.Logger

func loadAppConfig() (common.Config, error) {
	var conf common.Config
	if _, err := toml.DecodeFile("config/configuration.toml", &conf); err != nil {
		appLogger.Fatal("Couldn't load configuration: %s", err.Error())
		return common.Config{}, err
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
