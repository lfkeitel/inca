package grabber

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/dragonrider23/infrastructure-config-archive/interfaces"
)

func loadDeviceList(conf interfaces.Config) ([]host, error) {
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
			appLogger.Error("Error on line %d in device configuration", lineNum)
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

func loadDeviceTypes(conf interfaces.Config) ([]dtype, error) {
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
			appLogger.Error("Error on line %d in device type configuration", lineNum)
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

func grabConfigs(hosts []host, dtypes []dtype, dateSuffix string, conf interfaces.Config) error {
	var wg sync.WaitGroup
	ccg := newConnGroup(conf) // Used to enforce a maximum number of connections

	for _, host := range hosts {
		host := host
		for _, dtype := range dtypes {
			if host.dtype == dtype.deviceType && host.method == dtype.method {
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
				break
			}
		}
		ccg.wait()
	}

	wg.Wait()
	return nil
}

func getConfigFileName(host host, dateSuffix string, conf interfaces.Config) string {
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

func getArguments(argStr string, host host, filename string, conf interfaces.Config) []string {
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
	_, err := exec.Command("scripts/"+sfn, args...).Output()
	if err != nil {
		appLogger.Error(err.Error())
	}
	//stdOutLogger.Info(string(out))
	return nil
}
