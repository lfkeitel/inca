package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/dragonrider23/inca/database"
	"github.com/dragonrider23/inca/internal/common"
	"github.com/dragonrider23/inca/logger"
	"github.com/dragonrider23/inca/poller"
	"github.com/dragonrider23/inca/server"
)

var (
	configFile string
	appLogger  = logger.New("core")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		appLogger.Info("Shutting Down...")
		poller.Close()
		database.Close()
		os.Exit(1)
	}()

	flag.StringVar(&configFile, "c", "config/configuration.toml", "Configuration file")
	flag.Parse()

	if err := common.LoadAppConfig(configFile); err != nil {
		appLogger.Fatal(err.Error())
	}
	database.Prepare()
	server.Start()
	return
}
