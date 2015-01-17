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
			name:         splitLine[0],
			address:      splitLine[1],
			manufacturer: splitLine[2],
			proto:        splitLine[3],
		}

		hostList = append(hostList, device)
	}

	return hostList, nil
}

func grabConfigs(hosts []host, dateSuffix string, conf interfaces.Config) error {
	var wg sync.WaitGroup
	ccg := newConnGroup(conf) // Used to enforce a maximum number of connections

	for _, host := range hosts {
		host := host

		if host.manufacturer == "juniper" {
			if host.proto == "ssh" {
				fname := prepareTftpFile(host, dateSuffix, conf)
				wg.Add(1)
				ccg.add(1)
				go func() {
					defer func() {
						wg.Done()
						ccg.done()
					}()
					juniperExecute(host, fname, conf)
				}()
			}
			continue
		} else if host.manufacturer == "cisco" {
			if host.proto == "ssh" {
				fname := prepareTftpFile(host, dateSuffix, conf)
				wg.Add(1)
				ccg.add(1)
				go func() {
					defer func() {
						wg.Done()
						ccg.done()
					}()
					ciscoExecute(host, fname, conf)
				}()

			} else if host.proto == "telnet" {
				fname := prepareTftpFile(host, dateSuffix, conf)
				wg.Add(1)
				ccg.add(1)
				go func() {
					defer func() {
						wg.Done()
						ccg.done()
					}()
					ciscoExecute(host, fname, conf)
				}()

			} else {
				appLogger.Error("Protocol %s is not supported on Cisco", host.proto)
			}
		} else {
			appLogger.Error("Manufacturer %s is not supported", host.manufacturer)
		}
		ccg.wait()
	}

	wg.Wait()
	return nil
}

func prepareTftpFile(host host, dateSuffix string, conf interfaces.Config) string {
	var filename bytes.Buffer

	filename.WriteString(host.name)
	filename.WriteString("-")
	filename.WriteString(dateSuffix)
	filename.WriteString("-")
	filename.WriteString(host.address)
	filename.WriteString("-")
	filename.WriteString(host.manufacturer)
	filename.WriteString("-")
	filename.WriteString(host.proto)
	filename.WriteString(".conf")

	touch(conf.FullConfDir + "/" + filename.String())
	err := os.Chmod(conf.FullConfDir+"/"+filename.String(), 0777)
	if err != nil {
		appLogger.Error(err.Error())
	}

	return filename.String()
}

func ciscoExecute(host host, filename string, conf interfaces.Config) error {
	cmd := exec.Command("scripts/cisco-"+host.proto+"-config-grab.exp", conf.Tftphost, host.address, conf.RemotePassword, conf.RemoteUsername, filename, conf.EnablePassword)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		appLogger.Error(err.Error())
	}
	stdOutLogger.Info(out.String())
	return nil
}
