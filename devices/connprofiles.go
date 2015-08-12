package devices

import (
	db "github.com/dragonrider23/inca/database"
)

func GetConnProfiles() ([]ConnProfile, error) {
	rows, err := db.Conn.Query("SELECT * FROM conn_profiles")
	if err != nil {
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
			return nil, err
		}

		connProfileList = append(connProfileList, c)
	}

	return connProfileList, nil
}

func GetConnProfile(id string) (ConnProfile, error) {
	row := db.Conn.QueryRow(`SELECT * FROM conn_profiles
		WHERE profileid = ?`, id)

	var c ConnProfile

	err := row.Scan(
		&c.Profileid,
		&c.Name,
		&c.Protocol,
		&c.Username,
		&c.Password,
		&c.Enable,
	)

	return c, err
}

func CreateConnProfile(c ConnProfile) error {
	_, err := db.Conn.Exec(`INSERT INTO conn_profiles
		VALUES (null, ?, ?, ?, ?, ?)`,
		c.Name,
		c.Protocol,
		c.Username,
		c.Password,
		c.Enable)

	return err
}

func EditConnProfile(c ConnProfile) error {
	_, err := db.Conn.Exec(`UPDATE conn_profiles
		SET name = ?,
			protocol = ?,
			username = ?,
			password = ?,
			enable = ?
		WHERE profileid = ?`,
		c.Name,
		c.Protocol,
		c.Username,
		c.Password,
		c.Enable,
		c.Profileid)

	return err
}

func DeleteConnProfile(id string) error {
	_, err := db.Conn.Exec(`DELETE FROM conn_profiles WHERE profileid = ?`, id)
	return err
}
