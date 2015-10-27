package poller

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	fs "github.com/dragonrider23/utils/filesystem"
	un "github.com/dragonrider23/utils/net"
	us "github.com/dragonrider23/utils/sync"

	db "github.com/dragonrider23/inca/database"
	"github.com/dragonrider23/inca/devices"
	"github.com/dragonrider23/inca/internal/common"
)

var (
	configJobLG = us.NewLimitGroup(1)
)

func configJob(j Job, out chan<- *Response, e chan<- error) {
	defer func() {
		if r := recover(); r != nil {
			e <- errors.New("Paniced")
		}
	}()

	configJobLG.Wait()
	configJobLG.Add(1)
	defer configJobLG.Done()

	//startTime := time.Now()

	hosts := j.Data["devices"]
	connProfiles := j.Data["connProfiles"]
	grabConfigs(hosts.([]devices.Device), connProfiles.([]devices.ConnProfile))

	//endTime := time.Now()
	return
}

func grabConfigs(hosts []devices.Device, connProfiles []devices.ConnProfile) error {
	var wg sync.WaitGroup
	conf := common.Config
	// Used to enforce a maximum number of connections
	lg := us.NewLimitGroup(int32(conf.MaxSimultaneousConn))

	for _, host := range hosts {
		host := host
		match := false

		for _, connProfile := range connProfiles {
			if host.ConnProfile == connProfile.Profileid &&
				connProfile.Protocol != "none" &&
				!host.Custom {
				if supportedProtocol(host.Manufacturer, connProfile.Protocol) {
					script := getScriptFilename(host.Manufacturer, connProfile.Protocol)
					filename, timestamp := getConfigFileName(host)
					args := []string{
						host.Hostname,
						connProfile.Username,
						connProfile.Password,
						filepath.Join(conf.DataDir, filename),
						connProfile.Enable,
					}

					fmt.Printf("Getting config for %s @ %s\n", host.Name, host.Hostname)
					wg.Add(1)
					lg.Add(1)
					go func() {
						defer func() {
							wg.Done()
							lg.Done()
						}()
						lastError := scriptExecute(script, args)
						saveStatusInfoToDb(host, filename, timestamp, lastError)
					}()
					match = true
					break
				}
			}
		}

		if !match {
			// Device has an invalid connection profile
			// or doesn't support configuration grabs
			fmt.Printf("Pinging %s @ %s\n", host.Name, host.Hostname)
			wg.Add(1)
			lg.Add(1)
			go func() {
				defer func() {
					wg.Done()
					lg.Done()
				}()
				pingDevice(host)
			}()
		}

		lg.Wait()
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
	conf := common.Config
	dateTime := time.Now()

	filename.WriteString("configs/")
	filename.WriteString(host.Hostname)
	filename.WriteString("/")
	filename.WriteString(dateTime.Format("20060102"))
	filename.WriteString("/")
	filename.WriteString(strconv.FormatInt(dateTime.Unix(), 10))
	filename.WriteString(".conf")

	if err := fs.Touch(conf.DataDir + "/" + filename.String()); err != nil {
		fmt.Println(err.Error())
	}
	return filename.String(), dateTime.Unix()
}

func scriptExecute(sfn string, args []string) string {
	errMsg := ""
	out, err := exec.Command("poller/scripts/"+sfn, args...).Output()
	if err != nil {
		os.Remove(args[3])
		fmt.Println(string(out))

		switch err.Error() {
		case "exit status 64":
			errMsg = "Timeout Exceeded"
		case "exit status 65":
			errMsg = "SSH Login Failed"
		case "exit status 66":
			errMsg = "Telnet Login Failed"
		case "exit status 67":
			errMsg = "Enable Mode Failed"
		default:
			errMsg = "Unknown Error"
		}
	}
	return errMsg
}

func saveStatusInfoToDb(device devices.Device, filename string, timestamp int64, lastError string) error {
	conf := common.Config
	status := 0

	if _, err := os.Stat(filepath.Join(conf.DataDir, filename)); os.IsNotExist(err) {
		status = 2
	}

	if err := devices.UpdateDeviceStatus(device, status, timestamp, lastError); err != nil {
		return err
	}

	if status == 0 {
		_, err := db.Conn.Exec(`INSERT INTO configs
    		VALUES (null, ?, ?, ?, 1, "")`, device.Deviceid, timestamp, filename)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func pingDevice(d devices.Device) error {
	status := 0
	lastError := ""

	_, err := un.Ping(d.Hostname, false)
	if err != nil {
		fmt.Println(err.Error())
		status = 2
		lastError = "Host Unreachable"
	}

	return devices.UpdateDeviceStatus(d, status, time.Now().Unix(), lastError)
}
