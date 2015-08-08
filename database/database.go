package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dragonrider23/infrastructure-config-archive/common"
)

var conf common.Config
var Conn *sql.DB

func Prepare(config common.Config) error {
	conf = config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Address,
		conf.Database.Port,
		conf.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		errors.New("Failed to open database")
	}
	Conn = db
	return nil
}

func Close() {
	Conn.Close()
	return
}
