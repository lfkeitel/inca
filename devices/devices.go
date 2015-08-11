package devices

import (
	"github.com/dragonrider23/utils/slices"

	db "github.com/dragonrider23/infrastructure-config-archive/database"
)

func GetDevicesForConfigs() ([]Device, error) {
	return getDevices(true, 0)
}

func GetAllDevices() ([]Device, error) {
	return getDevices(false, 0)
}

func getDevices(disabled bool, id int) ([]Device, error) {
	var args []interface{}
	statement := `SELECT d.*, s.status, s.last_polled
		FROM devices AS d
		LEFT JOIN device_status AS s
		ON d.deviceid = s.deviceid`

	if disabled {
		statement += " WHERE d.disabled = 0"
	}

	if id > 0 {
		statement += " WHERE d.deviceid = ?"
		args = append(args, id)
	}

	rows, err := db.Conn.Query(statement, args...)
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
	d, err := getDevices(false, id)
	if err != nil {
		return Device{}, err
	}

	if len(d) != 1 {
		return Device{}, newError("No device found")
	}
	return d[0], err
}

func CreateDevice(d Device) error {
	d.Custom = !isSupportedDevice(d)

	_, err := db.Conn.Exec(`INSERT INTO devices
		VALUES (null, ?, ?, ?, ?, ?, ?, ?)`,
		d.Name,
		d.Hostname,
		d.ConnProfile,
		d.Manufacturer,
		d.Model,
		d.Custom,
		d.Disabled)

	return err
}

func EditDevice(d Device) error {
	d.Custom = !isSupportedDevice(d)

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

func DeleteDevices(id []int) error {
	statement := "DELETE FROM devices WHERE 0"

	for range id {
		statement += " OR deviceid = ?"
	}

	ids := convertIntSliceToInterface(id)
	_, err := db.Conn.Exec(statement, ids...)
	return err
}

// A supported device here is if Inca will get the configuration off the device
func isSupportedDevice(d Device) bool {
	ma, ok := supportedDeviceTypes[d.Manufacturer]
	return ok && slices.StringInSlice(d.Model, ma)
}

func convertIntSliceToInterface(s []int) []interface{} {
	is := make([]interface{}, len(s))
	for i, d := range s {
		is[i] = d
	}
	return is
}
