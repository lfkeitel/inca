package grabber

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lfkeitel/inca/src/common"
	"github.com/lfkeitel/inca/src/diff"
)

func loadDeviceList(conf *common.Config) ([]host, error) {
	listFile, err := os.Open(conf.Paths.DeviceList)
	if err != nil {
		return nil, err
	}
	defer listFile.Close()

	scanner := bufio.NewScanner(listFile)
	scanner.Split(bufio.ScanLines)
	var hostList []host
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
			common.UserLogWarning(logText)
			continue
		}

		device := host{
			Name:    splitLine[0],
			Address: splitLine[1],
			Dtype:   splitLine[2],
			Method:  splitLine[3],
		}

		hostList = append(hostList, device)
	}

	return hostList, nil
}

// CheckDeviceList parses string s as if it was a device/type list and checks for errors
func CheckDeviceList(s string) error {
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		if len(line) == 0 || line[0] == '#' || line[0] == ' ' {
			continue
		}

		parsedLine := strings.Split(line, "::")
		if len(parsedLine) != 4 {
			return fmt.Errorf("Error on line %d. Expected 4 fields, got %d.\\n'%s'", i+1, len(parsedLine), lines[i])
		}
	}
	return nil
}

func loadDeviceTypes(conf *common.Config) ([]dtype, error) {
	typeFile, err := os.Open(conf.Paths.DeviceTypes)
	if err != nil {
		return nil, err
	}
	defer typeFile.Close()

	scanner := bufio.NewScanner(typeFile)
	scanner.Split(bufio.ScanLines)
	var dtypeList []dtype
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
			common.UserLogWarning(logText)
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

func grabConfigs(hosts []host, dtypes []dtype, dateSuffix string, conf *common.Config) error {
	var wg sync.WaitGroup
	ccg := newConnGroup(conf) // Used to enforce a maximum number of connections
	fname := dateSuffix + ".conf"

	for _, host := range hosts {
		host := host
		match := false
		for _, dtype := range dtypes {
			if host.Dtype == dtype.deviceType && (dtype.method == "*" || host.Method == dtype.method) {
				hostdir := filepath.Join(conf.Paths.ConfDir, fmt.Sprintf("%s-%s", host.Name, host.Address))
				hostfname := filepath.Join(hostdir, fname)
				metadatafile := filepath.Join(hostdir, "_metadata.json")
				args := getArguments(dtype.args, host, hostfname, conf)

				if !common.FileExists(hostdir) {
					if err := os.MkdirAll(hostdir, 0750); err != nil {
						common.UserLogError(err.Error())
						break
					}
				}

				if !common.FileExists(metadatafile) {
					hostmetafile := metadatafile
					d, _ := json.Marshal(host)
					if err := ioutil.WriteFile(hostmetafile, d, 0640); err != nil {
						common.UserLogError(err.Error())
						break
					}
				}

				wg.Add(1)
				ccg.add(1)
				go func() {
					defer func() {
						appLogger.Debugf("Done with %s", host.Name)
						wg.Done()
						ccg.done()
					}()

					// Get latest config filename before new one is created
					dircontents, _ := ioutil.ReadDir(hostdir)
					latestConfig := ""
					for _, f := range dircontents {
						if f.Name() != "_metadata.json" {
							latestConfig = f.Name()
						}
					}

					// Get new config
					if err := scriptExecute(dtype.scriptfile, args); err != nil {
						common.UserLogError("Failed getting config from %s (%s)", host.Name, host.Address)
						os.Remove(hostfname)
						return
					}

					// If an old config existed, check if it's different from the new one
					if latestConfig != "" {
						same, err := diff.SameFileContents(hostfname, filepath.Join(hostdir, latestConfig))
						if err != nil {
							appLogger.Error(err.Error())
							return
						}

						if same {
							appLogger.Debug(
								"(%s) Last config and this config are the same, deleting file",
								host.Name,
							)
							os.Remove(hostfname)
						}
					}
				}()
				match = true
				break
			}
		}

		if !match {
			logText := fmt.Sprintf("Device type '%s' using method '%s' wasn't found.", host.Dtype, host.Method)
			appLogger.Warning(logText)
			common.UserLogWarning(logText)
			finishedDevices++
		}
		appLogger.Debug("Waiting for available slot")
		ccg.wait()
	}

	appLogger.Debug("Waiting for all devices")
	wg.Wait()
	appLogger.Debug("All devices finished")
	return nil
}

func cleanUpHostDirs(hosts []host) {
	for _, h := range hosts {
		hostdir := filepath.Join(conf.Paths.ConfDir, fmt.Sprintf("%s-%s", h.Name, h.Address))
		dirlist, _ := ioutil.ReadDir(hostdir)

		if len(dirlist) > conf.KeepLimit+1 { // +1 for metadata file
			dirlist := filterStrings(
				dirlistToFilenames(dirlist),
				func(s string) bool { return s != "_metadata.json" },
			)

			for _, f := range dirlist[:len(dirlist)-conf.KeepLimit] {
				os.Remove(filepath.Join(hostdir, f))
			}
		}
	}
}

func getConfigFileName(host host, dateSuffix string, conf *common.Config) string {
	filename := fmt.Sprintf("%s-%s-%s-%s-%s.conf", host.Name, dateSuffix, host.Address, host.Dtype, host.Method)
	confPath := filepath.Join(conf.Paths.ConfDir, filename)

	touch(confPath)
	return confPath
}

func getArguments(argStr string, host host, filename string, conf *common.Config) []string {
	args := strings.Split(argStr, ",")
	argList := make([]string, len(args))
	for i, a := range args {
		switch a {
		case "$address":
			argList[i] = host.Address
			break
		case "$username":
			argList[i] = conf.Credentials.RemoteUsername
			break
		case "$password":
			argList[i] = conf.Credentials.RemotePassword
			break
		case "$logfile":
			argList[i] = filename
			break
		case "$enablepw":
			argList[i] = conf.Credentials.EnablePassword
			break
		}
	}
	return argList
}

func scriptExecute(sfn string, args []string) error {
	out, err := exec.Command(filepath.Join(conf.Paths.ScriptDir, sfn), args...).Output()
	if err != nil {
		appLogger.Error(err.Error())
		appLogger.Error(string(out))
		return err
	}
	return nil
}
