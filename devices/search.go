package devices

import (
	"net"

	db "github.com/dragonrider23/inca/database"
)

// Search searchs the database for devices using query q. If q is a valid
// IP address then it will search the hostname, otherwise it will
// search the device names.
func Search(q string) ([]Device, error) {
	statement := `SELECT d.*, s.status, s.last_polled, s.last_error
		FROM devices AS d
		LEFT JOIN device_status AS s
		ON d.deviceid = s.deviceid`

	if isIPAddress(q) {
		statement += " WHERE d.hostname = ?"
	} else {
		statement += " WHERE d.name LIKE ?"
		q = "%" + q + "%"
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
			&d.Status.LastError,
		)

		if err != nil {
			return nil, err
		}

		deviceList = append(deviceList, d)
	}

	return deviceList, nil
}

func isIPAddress(a string) bool {
	ip := net.ParseIP(a)
	return ip != nil
}
