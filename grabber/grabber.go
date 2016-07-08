package grabber

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/lfkeitel/inca/comm"
)

func loadDeviceList(conf comm.Config) ([]host, error) {
	listFile, err := os.Open(conf.DeviceListFile)
	if err != nil {
		return nil, err
	}
	defer listFile.Close()

	scanner := bufio.NewScanner(listFile)
	scanner.Split(bufio.ScanLines)
	hostList := make([]host, 0)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if len(line) < 1 || line[0] == '#' || line[0] == ' ' {
			continue
		}

		splitLine := strings.Split(line, "::")

		if len(splitLine) != 4 {
			logText := fmt.Sprintf("Error on line %d in device configuration", lineNum)
			appLogger.Warning(logText)
			comm.UserLogWarning(logText)
			continue
		}

		device := host{
			name:    splitLine[0],
			address: splitLine[1],
			dtype:   splitLine[2],
			method:  splitLine[3],
		}

		hostList = append(hostList, device)
	}

	return hostList, nil
}

// Parses string s as if it was a device [type] list and checks for errors
func CheckDeviceList(s string) error {
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		if len(line) < 1 || line[0] == '#' || line[0] == ' ' {
			continue
		}

		parsedLine := strings.Split(line, "::")
		if len(parsedLine) != 4 {
			return fmt.Errorf("Error on line %d. Expected 4 fields, got %d.\\n'%s'", i+1, len(parsedLine), lines[i])
		}
	}
	return nil
}

func loadDeviceTypes(conf comm.Config) ([]dtype, error) {
	typeFile, err := os.Open(conf.DeviceTypeFile)
	if err != nil {
		return nil, err
	}
	defer typeFile.Close()

	scanner := bufio.NewScanner(typeFile)
	scanner.Split(bufio.ScanLines)
	dtypeList := make([]dtype, 0)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if len(line) < 1 || line[0] == '#' || line[0] == ' ' {
			continue
		}

		splitLine := strings.Split(line, "::")

		if len(splitLine) != 4 {
			logText := fmt.Sprintf("Error on line %d in device type configuration", lineNum)
			appLogger.Warning(logText)
			comm.UserLogWarning(logText)
			continue
		}

		typedef := dtype{
			deviceType: splitLine[0],
			method:     splitLine[1],
			scriptfile: splitLine[2],
			args:       splitLine[3],
		}

		dtypeList = append(dtypeList, typedef)
	}

	return dtypeList, nil
}

func grabConfigs(hosts []host, dtypes []dtype, dateSuffix string, conf comm.Config) error {
	var wg sync.WaitGroup
	ccg := newConnGroup(conf) // Used to enforce a maximum number of connections

	for _, host := range hosts {
		host := host
		match := false
		for _, dtype := range dtypes {
			if host.dtype == dtype.deviceType && (dtype.method == "*" || host.method == dtype.method) {
				fname := getConfigFileName(host, dateSuffix, conf)
				args := getArguments(dtype.args, host, fname, conf)

				wg.Add(1)
				ccg.add(1)
				go func() {
					defer func() {
						wg.Done()
						ccg.done()
					}()
					scriptExecute(dtype.scriptfile, args)
				}()
				match = true
				break
			}
		}

		if !match {
			logText := fmt.Sprintf("Device type '%s' using method '%s' wasn't found.", host.dtype, host.method)
			appLogger.Warning(logText)
			comm.UserLogWarning(logText)
			finishedDevices++
		}
		ccg.wait()
	}

	wg.Wait()
	return nil
}

func getConfigFileName(host host, dateSuffix string, conf comm.Config) string {
	var filename bytes.Buffer

	filename.WriteString(conf.FullConfDir)
	filename.WriteString("/")
	filename.WriteString(host.name)
	filename.WriteString("-")
	filename.WriteString(dateSuffix)
	filename.WriteString("-")
	filename.WriteString(host.address)
	filename.WriteString("-")
	filename.WriteString(host.dtype)
	filename.WriteString("-")
	filename.WriteString(host.method)
	filename.WriteString(".conf")

	touch(conf.FullConfDir + "/" + filename.String())

	return filename.String()
}

func getArguments(argStr string, host host, filename string, conf comm.Config) []string {
	args := strings.Split(argStr, ",")
	argList := make([]string, len(args))
	for i, a := range args {
		switch a {
		case "$address":
			argList[i] = host.address
			break
		case "$username":
			argList[i] = conf.RemoteUsername
			break
		case "$password":
			argList[i] = conf.RemotePassword
			break
		case "$logfile":
			argList[i] = filename
			break
		case "$enablepw":
			argList[i] = conf.EnablePassword
			break
		}
	}
	return argList
}

func scriptExecute(sfn string, args []string) error {
	out, err := exec.Command("scripts/"+sfn, args...).Output()
	if err != nil {
		appLogger.Error(err.Error())
		appLogger.Error(string(out))
		comm.UserLogError("Failed getting config from %s", args[0])
	}
	return nil
}
