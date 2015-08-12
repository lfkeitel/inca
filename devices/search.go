package devices

import (
	"regexp"

	db "github.com/dragonrider23/inca/database"
)

const (
	ipRegex = `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
)

func Search(q string) ([]Device, error) {
	var args []string
	statement := `SELECT d.*, s.status, s.last_polled
		FROM devices AS d
		LEFT JOIN device_status AS s
		ON d.deviceid = s.deviceid`

	if isIPAddress(q) {
		statement += " WHERE d.hostname = ?"
		args = append(args, q)
	} else {
		statement += " WHERE d.name LIKE ?"
		q = "%" + q + "%"
		args = append(args, q)
	}

	rows, err := db.Conn.Query(statement, q)
	if err != nil {
		return nil, err
	}

	var deviceList []Device

	for rows.Next() {
		var d Device
		err = rows.Scan(
			&d.Deviceid,
			&d.Name,
			&d.Hostname,
			&d.ConnProfile,
			&d.Manufacturer,
			&d.Model,
			&d.Custom,
			&d.Disabled,
			&d.ParseConfig,
			&d.Status.Status,
			&d.Status.LastPolled,
		)

		if err != nil {
			return nil, err
		}

		deviceList = append(deviceList, d)
	}

	return deviceList, nil
}

func isIPAddress(a string) bool {
	errRegEx, rerr := regexp.Compile(ipRegex)
	if rerr != nil {
		return false
	}

	line := errRegEx.FindStringSubmatch(a)
	if line == nil {
		return false
	}

	return true
}
