package devices

import (
	"strings"

	"github.com/dragonrider23/utils/slices"

	db "github.com/dragonrider23/inca/database"
)

func GetDevicesForConfigGrab(ids []int) ([]Device, error) {
	var args []interface{}
	statement := `SELECT d.*, s.status, s.last_polled, s.last_error
		FROM devices AS d
		LEFT JOIN device_status AS s
		ON d.deviceid = s.deviceid
		WHERE 0`

	if ids != nil {
		for _, i := range ids {
			statement += " OR d.deviceid = ?"
			args = append(args, i)
		}
	}

	return getDevicesFromQuery(statement, args)
}

func GetAllDevices() ([]Device, error) {
	return getDevices(false, 0, -1)
}

func getDevices(disabled bool, id int, status int) ([]Device, error) {
	var args []interface{}
	statement := `SELECT d.*, s.status, s.last_polled, s.last_error
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

	if status > -1 {
		statement += " WHERE s.status != ?"
		args = append(args, status)
	}

	return getDevicesFromQuery(statement, args)
}

func getDevicesFromQuery(statement string, args []interface{}) ([]Device, error) {
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

func GetDevice(id int) (Device, error) {
	d, err := getDevices(false, id, -1)
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
	d.ParseConfig = isParsableDevice(d)

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	r, err := tx.Exec(`INSERT INTO devices
		VALUES (null, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.Name,
		d.Hostname,
		d.ConnProfile,
		d.Manufacturer,
		d.Model,
		d.Custom,
		d.Disabled,
		d.ParseConfig)
	if err != nil {
		tx.Rollback()
		return err
	}

	dID, err := r.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`INSERT INTO device_status
		VALUES (null, ?, 1, 0, "")`, dID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func EditDevice(d Device) error {
	d.Custom = !isSupportedDevice(d)
	d.ParseConfig = isParsableDevice(d)

	_, err := db.Conn.Exec(`UPDATE devices
		SET name = ?,
			hostname = ?,
			conn_profile = ?,
			manufacturer = ?,
			model = ?,
			custom = ?,
			disabled = ?,
			parse_config = ?
		WHERE deviceid = ?`,
		d.Name,
		d.Hostname,
		d.ConnProfile,
		d.Manufacturer,
		d.Model,
		d.Custom,
		d.Disabled,
		d.ParseConfig,
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
	man := strings.ToLower(d.Manufacturer)
	ma, ok := supportedDeviceTypes[man]
	return ok && slices.StringInSlice(d.Model, ma)
}

// A parsable device here is if Inca can parse the config for more advanced management
func isParsableDevice(d Device) bool {
	man := strings.ToLower(d.Manufacturer)
	ma, ok := parsableDeviceTypes[man]
	return ok && slices.StringInSlice(d.Model, ma)
}

func GetDeviceStats() (DeviceStatus, error) {
	devices, err := getDevices(false, 0, 0)
	if err != nil {
		return DeviceStatus{}, err
	}
	devLength := len(devices)
	var total int
	var unknownTotal int

	r := db.Conn.QueryRow(`SELECT COUNT(statusid)
		FROM device_status`)
	err = r.Scan(&total)
	if err != nil {
		return DeviceStatus{}, err
	}

	r = db.Conn.QueryRow(`SELECT COUNT(statusid)
		FROM device_status
		WHERE status = 1`)
	err = r.Scan(&unknownTotal)
	if err != nil {
		return DeviceStatus{}, err
	}

	d := DeviceStatus{
		Total:       total,
		Down:        devLength - unknownTotal,
		Up:          total - devLength,
		Unknown:     unknownTotal,
		DownDevices: devices,
	}

	return d, nil
}

func UpdateDeviceStatus(d Device, s int, t int64, le string) error {
	_, err := db.Conn.Exec(`UPDATE device_status
        SET status = ?,
            last_polled = ?,
            last_error = ?
        WHERE deviceid = ?`,
		s,
		t,
		le,
		d.Deviceid)
	if err != nil {
		return err
	}
	return nil
}
