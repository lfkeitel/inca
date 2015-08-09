package devices

import (
	"github.com/dragonrider23/utils/slices"

	db "github.com/dragonrider23/infrastructure-config-archive/database"
)

func GetDevicesForConfigs() ([]Device, error) {
	rows, err := db.Conn.Query(`SELECT d.*, s.status, s.last_polled
		FROM devices AS d
		LEFT JOIN device_status AS s
		ON d.deviceid = s.deviceid
		WHERE d.disabled = 0 AND d.custom = 0`)
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

func GetAllDevices() ([]Device, error) {
	rows, err := db.Conn.Query(`SELECT d.*, s.status, s.last_polled
		FROM devices AS d
		LEFT JOIN device_status AS s
		ON d.deviceid = s.deviceid`)
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

func GetDevice(id int) (Device, error) {
	row := db.Conn.QueryRow(`SELECT d.*, s.status, s.last_polled
		FROM devices AS d
		LEFT JOIN device_status AS s
		ON d.deviceid = s.deviceid
		WHERE d.deviceid = ?`, id)

	var d Device

	err := row.Scan(
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

	return d, err
}

func CreateDevice(d Device) error {
	d.Custom = !isSupportedDevice(d)

	_, err := db.Conn.Exec(`INSERT INTO devices
		VALUES (null, ?, ?, ?, ?, ?, ?)`,
		d.Name,
		d.Hostname,
		d.ConnProfile,
		d.Manufacturer,
		d.Model,
		d.Custom)

	return err
}

func EditDevice(d Device) error {
	_, err := db.Conn.Exec(`UPDATE devices
		SET name = ?,
			hostname = ?,
			conn_profile = ?,
			manufacturer = ?,
			model = ?,
			custom = ?,
			disabled = ?
		WHERE deviceid = ?`,
		d.Name,
		d.Hostname,
		d.ConnProfile,
		d.Manufacturer,
		d.Model,
		d.Custom,
		d.Disabled,
		d.Deviceid)

	return err
}

func DeleteDevice(id string) error {
	_, err := db.Conn.Exec(`DELETE FROM devices WHERE deviceid = ?`, id)
	return err
}

func isSupportedDevice(d Device) bool {
	ma, ok := supportedDeviceTypes[d.Manufacturer]
	return ok && slices.StringInSlice(d.Model, ma)
}
