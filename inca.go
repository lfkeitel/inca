package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"syscall"
	//"time"

	"github.com/naoina/toml"

	"github.com/dragonrider23/inca/common"
	"github.com/dragonrider23/inca/configs"
	"github.com/dragonrider23/inca/database"
	"github.com/dragonrider23/inca/logger"
	"github.com/dragonrider23/inca/server"
)

const (
	configFile = "config/configuration.toml"
)

var (
	appLogger *logger.Logger
)

func loadAppConfig(fn string) (common.Config, error) {
	var conf common.Config

	f, err := ioutil.ReadFile(fn)
	if err != nil {
		return common.Config{}, err
	}

	if err = toml.Unmarshal(f, &conf); err != nil {
		// Attempt to print a meaningful error message
		errRegEx, rerr := regexp.Compile(`^toml:.*?line (\d+):`)
		if rerr != nil {
			appLogger.Fatal("Invalid configuration. %s", err.Error())
			return common.Config{}, err
		}

		line := errRegEx.FindStringSubmatch(err.Error())
		if line == nil {
			appLogger.Fatal("Invalid configuration. %s", err.Error())
			return common.Config{}, err
		}

		appLogger.Fatal("Invalid configuration. Check line %s.", line[1])
		return common.Config{}, err
	}
	return conf, nil
}

func init() {
	appLogger = logger.New("core")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		appLogger.Info("Shutting Down...")
		configs.Stop()
		database.Close()
		os.Exit(1)
	}()

	conf, _ := loadAppConfig(configFile)
	database.Prepare(&conf)

	if err := configs.Prepare(&conf); err != nil {
		appLogger.Fatal("Faild to start poller: %s", err.Error())
	}

	server.Start(&conf)
	return
}