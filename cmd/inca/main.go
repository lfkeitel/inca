package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/lfkeitel/inca/src/common"
	"github.com/lfkeitel/inca/src/grabber"
	"github.com/lfkeitel/inca/src/server"
	"github.com/lfkeitel/inca/src/targz"
	"github.com/lfkeitel/verbose"
)

var (
	appLogger *verbose.Logger

	showVersion bool
	configFile  string

	version   = ""
	buildTime = ""
	builder   = ""
	goversion = ""
)

func init() {
	flag.BoolVar(&showVersion, "v", false, "Print version information")
	flag.StringVar(&configFile, "c", "config/configuration.toml", "Configuration file")
}

func main() {
	flag.Parse()

	if showVersion {
		displayVersionInfo()
		return
	}

	conf, _ := common.LoadConfig(configFile, appLogger)

	setupLogger(conf.Paths.LogDir)

	common.InitUserLog(conf.Paths.LogDir)
	grabber.Init(conf)
	server.StartServer(conf)
}

func setupLogger(logdir string) {
	appLogger = verbose.New("app")

	fileLogger, err := verbose.NewFileHandler(filepath.Join(logdir, "app.log"))
	if err != nil {
		panic("Failed to open logging directory")
	}

	appLogger.AddHandler("file", fileLogger)
	appLogger.AddHandler("stdout", verbose.NewStdoutHandler(true))
}

func setupTarLogger(logdir string) {
	appLogger := verbose.New("tarGz-log")

	fileLogger, err := verbose.NewFileHandler("logs/tar.log")
	if err != nil {
		panic("Failed to open logging directory")
	}

	appLogger.AddHandler("file", fileLogger)
	appLogger.AddHandler("stdout", verbose.NewStdoutHandler(true))
	targz.SetLogger(appLogger)
}

func displayVersionInfo() {
	fmt.Printf(`INCA - (C) 2017 University of Southern Indiana - Lee Keitel

Version:     %s
Built:       %s
Compiled by: %s
Go version:  %s
`, version, buildTime, builder, goversion)
}
