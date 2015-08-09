package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"

	"github.com/dragonrider23/go-logger"
	"github.com/dragonrider23/infrastructure-config-archive/common"
	"github.com/dragonrider23/infrastructure-config-archive/configs"
	"github.com/dragonrider23/infrastructure-config-archive/database"
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		appLogger.Info("Shutting Down...")
		database.Close()
		os.Exit(1)
	}()

	conf, _ := loadAppConfig()
	configs.LoadConfig(conf)
	database.Prepare(conf)
	server.Start(conf)
	return
}
