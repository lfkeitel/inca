package common

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/lfkeitel/verbose"
)

type Config struct {
	MaxSimultaneousConn int
	Credentials         credentialsConf
	Paths               pathsConf
	Server              serverConf
	Hooks               hookConf
}

type credentialsConf struct {
	RemoteUsername string
	RemotePassword string
	EnablePassword string
}

type serverConf struct {
	BindAddress string
	BindPort    int
}

type pathsConf struct {
	DeviceList  string
	DeviceTypes string
	ConfDir     string
	ScriptDir   string
	ArchiveDir  string
	LogDir      string
}

type hookConf struct {
	PreScript  string
	PostScript string
}

func InitUserLog(logdir string) {
	userLog := verbose.New("endUserLog")

	fileLogger, err := verbose.NewFileHandler(filepath.Join(logdir, "endUser.log"))
	if err != nil {
		panic("Failed to open logging directory")
	}

	userLog.AddHandler("file", fileLogger)
}

type appLogger interface {
	Fatalf(string, ...interface{})
}

func LoadConfig(path string, logger appLogger) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		logger.Fatalf("Couldn't load configuration: %s", err.Error())
		return nil, err
	}

	if err := setDefaults(&conf); err != nil {
		return nil, err
	}

	return &conf, makeDirectories(&conf)
}

func setDefaults(c *Config) error {
	c.MaxSimultaneousConn = intOrDefault(c.MaxSimultaneousConn, 1000)

	c.Paths.DeviceList = stringOrDefault(c.Paths.DeviceList, "config/device-definitions.conf")
	c.Paths.DeviceTypes = stringOrDefault(c.Paths.DeviceTypes, "config/device-types.conf")
	c.Paths.ConfDir = stringOrDefault(c.Paths.ConfDir, "latest")
	c.Paths.ArchiveDir = stringOrDefault(c.Paths.ArchiveDir, "archive")
	c.Paths.ScriptDir = stringOrDefault(c.Paths.ScriptDir, "scripts")
	c.Paths.LogDir = stringOrDefault(c.Paths.LogDir, "logs")

	c.Server.BindPort = intOrDefault(c.Server.BindPort, 8080)
	return nil
}

func makeDirectories(c *Config) error {
	if !FileExists(c.Paths.ConfDir) {
		if err := os.MkdirAll(c.Paths.ConfDir, 0755); err != nil {
			return err
		}
	}

	if !FileExists(c.Paths.ArchiveDir) {
		if err := os.MkdirAll(c.Paths.ArchiveDir, 0755); err != nil {
			return err
		}
	}

	if !FileExists(c.Paths.ScriptDir) {
		if err := os.MkdirAll(c.Paths.ScriptDir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func stringOrDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func intOrDefault(s, def int) int {
	if s == 0 {
		return def
	}
	return s
}

func ReverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func UserLogInfo(format string, v ...interface{}) {
	verbose.Get("endUserLog").Infof(format, v...)
}

func UserLogWarning(format string, v ...interface{}) {
	verbose.Get("endUserLog").Warningf(format, v...)
}

func UserLogError(format string, v ...interface{}) {
	verbose.Get("endUserLog").Errorf(format, v...)
}

func UserLogFatal(format string, v ...interface{}) {
	verbose.Get("endUserLog").Fatalf(format, v...)
}
