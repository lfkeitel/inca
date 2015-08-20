package common

import (
	"errors"
	"io/ioutil"
	"regexp"

	"github.com/naoina/toml"
)

// Config is the application-wide configuration
var Config *Configuration

// LoadAppConfig reads and parses the file fn as the application configuration
func LoadAppConfig(fn string) error {
	var c Configuration
	f, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	if err = toml.Unmarshal(f, &c); err != nil {
		// Attempt to print a meaningful error message
		errRegEx, rerr := regexp.Compile(`^toml:.*?line (\d+):`)
		if rerr != nil {
			return errors.New("Invalid configuration. " + err.Error())
		}

		line := errRegEx.FindStringSubmatch(err.Error())
		if line == nil {
			return errors.New("Invalid configuration. " + err.Error())
		}

		return errors.New("Invalid configuration. Check line " + line[1])
	}
	Config = &c
	return nil
}
