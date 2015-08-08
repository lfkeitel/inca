package configs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/dragonrider23/utils/channels"
	fs "github.com/dragonrider23/utils/filesystem"

	"github.com/dragonrider23/infrastructure-config-archive/common"
	db "github.com/dragonrider23/infrastructure-config-archive/database"
	"github.com/dragonrider23/infrastructure-config-archive/devices"
)

func grabConfigs(hosts []devices.Device, connProfiles []devices.ConnProfile) error {
	var wg sync.WaitGroup
	// Used to enforce a maximum number of connections
	mcg := channels.NewMaxChanGroup(conf.MaxSimultaneousConn)

	for _, host := range hosts {
		host := host
		match := false

		for _, connProfile := range connProfiles {
			if host.ConnProfile == connProfile.Profileid {
				if supportedProtocol(host.Manufacturer, connProfile.Protocol) {
					script := getScriptFilename(host.Manufacturer, connProfile.Protocol)
					filename, timestamp := getConfigFileName(host)
					args := []string{
						host.Hostname,
						connProfile.Username,
						connProfile.Password,
						conf.DataDir + "/" + filename,
						connProfile.Enable,
					}

					wg.Add(1)
					mcg.Add(1)
					go func() {
						defer func() {
							wg.Done()
							finishedDevices++
							mcg.Done()
						}()
						scriptExecute(script, args)
						saveConfigInfoToDb(host, filename, timestamp)
					}()
					match = true
					break
				}
			}
		}

		if !match {
			logText := fmt.Sprintf("Misconfigured device or unsupported connection profile.")
			appLogger.Warning(logText)
			common.UserLogWarning(logText)
			finishedDevices++
		}

		mcg.Wait()
	}

	wg.Wait()
	return nil
}

func supportedProtocol(brand string, proto string) bool {
	if brand == "Cisco" {
		if proto == "ssh" || proto == "telnet" {
			return true
		}
	} else if brand == "Juniper" && proto == "ssh" {
		return true
	}

	return false
}

func getScriptFilename(brand string, proto string) string {
	brand = strings.ToLower(brand)
	return brand + "-" + proto + "-config-grab.exp"
}

func getConfigFileName(host devices.Device) (string, int64) {
	// Generate file name with the hierarchy
	// data/configs/$hostname/$date/$time.conf
	var filename bytes.Buffer
	dateTime := time.Now()

	filename.WriteString("configs/")
	filename.WriteString(host.Hostname)
	filename.WriteString("/")
	filename.WriteString(dateTime.Format("20060102"))
	filename.WriteString("/")
	filename.WriteString(dateTime.Format("15:04:05"))
	filename.WriteString(".conf")

	if err := fs.Touch(conf.DataDir + "/" + filename.String()); err != nil {
		appLogger.Error(err.Error())
	}
	return filename.String(), dateTime.Unix()
}

func scriptExecute(sfn string, args []string) error {
	out, err := exec.Command("configs/scripts/"+sfn, args...).Output()
	if err != nil {
		appLogger.Error(err.Error())
		appLogger.Error(string(out))
	}
	return nil
}

func saveConfigInfoToDb(device devices.Device, filename string, timestamp int64) error {
	_, err := db.Conn.Exec(`INSERT INTO configs
		VALUES (null, ?, ?, ?, 1, "")`, device.Deviceid, timestamp, filename)
	if err != nil {
		appLogger.Error(err.Error())
		return err
	}

	return nil
}
