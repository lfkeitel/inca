package grabber

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/dragonrider23/infrastructure-config-archive/common"
)

func loadDeviceList() ([]device, error) {
	// Get device information from database
	return nil, nil
}

func grabConfigs(hosts []device, dateSuffix string) error {
	var wg sync.WaitGroup
	mcg := newMaxChanGroup(conf.MaxSimultaneousConn) // Used to enforce a maximum number of connections

	for _, host := range hosts {
		host := host
		match := false
		// Check if connection profile exists
		// Check if connection profile matches supported protocols
		// Get the data filename
		// Get the arguments for the expect script
		// Run the script in a concurrent function
		//
		// for _, dtype := range dtypes {
		// 	if host.dtype == dtype.deviceType && (dtype.method == "*" || host.method == dtype.method) {
		// 		fname := getConfigFileName(host, dateSuffix)
		// 		args := []string{
		// 			host.address,
		// 			conf.RemoteUsername,
		// 			conf.RemotePassword,
		// 			fname,
		// 			conf.EnablePassword,
		// 		}
		//
		// 		wg.Add(1)
		// 		mcg.add(1)
		// 		go func() {
		// 			defer func() {
		// 				wg.Done()
		// 				mcg.done()
		// 			}()
		// 			scriptExecute(dtype.scriptfile, args)
		// 		}()
		// 		match = true
		// 		break
		// 	}
		// }
		//
		// if !match {
		// 	logText := fmt.Sprintf("Device type '%s' using method '%s' wasn't found.", host.dtype, host.method)
		// 	appLogger.Warning(logText)
		// 	common.UserLogWarning(logText)
		// 	finishedDevices++
		// }
		// mcg.wait()
	}

	wg.Wait()
	return nil
}

func getConfigFileName(host device, dateSuffix string) string {
	// Generate file name with the hierarchy
	// data/configs/$hostname/$date/$time.conf
	var filename bytes.Buffer

	filename.WriteString("data/configs/")
	filename.WriteString(host.address)
	filename.WriteString("/")
	filename.WriteString(dateSuffix)
	filename.WriteString("/")
	filename.WriteString("time") // replace with a timestamp
	filename.WriteString(".conf")

	touch(filename.String())
	return filename.String()
}

func scriptExecute(sfn string, args []string) error {
	out, err := exec.Command("scripts/"+sfn, args...).Output()
	if err != nil {
		appLogger.Error(err.Error())
		appLogger.Error(string(out))
	}
	return nil
}
