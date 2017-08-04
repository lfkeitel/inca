package main

import (
	"flag"
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/lfkeitel/inca/common"
	"github.com/lfkeitel/inca/grabber"
	"github.com/lfkeitel/inca/server"
	"github.com/lfkeitel/verbose"
)

var (
	appLogger *verbose.Logger

	showVersion bool

	version   = ""
	buildTime = ""
	builder   = ""
	goversion = ""
)

func init() {
	flag.BoolVar(&showVersion, "v", false, "Print version information")

	appLogger = verbose.New("app")

	fileLogger, err := verbose.NewFileHandler("logs/app.log")
	if err != nil {
		panic("Failed to open logging directory")
	}

	appLogger.AddHandler("file", fileLogger)
}

func main() {
	flag.Parse()

	if showVersion {
		displayVersionInfo()
		return
	}

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

func displayVersionInfo() {
	fmt.Printf(`INCA - (C) 2017 University of Southern Indiana - Lee Keitel

Version:     %s
Built:       %s
Compiled by: %s
Go version:  %s
`, version, buildTime, builder, goversion)
}
