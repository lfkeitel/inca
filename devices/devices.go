package devices

import (
	db "github.com/dragonrider23/infrastructure-config-archive/database"
)

func GetAllDevices() ([]Device, error) {
	rows, err := db.Conn.Query(`SELECT d.*, s.status, s.last_polled
		FROM devices AS d
		LEFT JOIN device_status AS s
		ON d.deviceid = s.deviceid
		WHERE d.disabled = 0 AND d.custom = 0`)
	if err != nil {
		// Log error
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
			&d.Status.Status,
			&d.Status.LastPolled,
		)

		if err != nil {
			// Log error
			return nil, err
		}

		deviceList = append(deviceList, d)
	}

	return deviceList, nil
}

func GetConnProfiles() ([]ConnProfile, error) {
	rows, err := db.Conn.Query("SELECT * FROM conn_profiles")
	if err != nil {
		// Log error
		return nil, err
	}

	var connProfileList []ConnProfile

	for rows.Next() {
		var c ConnProfile
		err = rows.Scan(
			&c.Profileid,
			&c.Name,
			&c.Protocol,
			&c.Username,
			&c.Password,
			&c.Enable,
		)

		if err != nil {
			// Log error
			return nil, err
		}

		connProfileList = append(connProfileList, c)
	}

	return connProfileList, nil
}
