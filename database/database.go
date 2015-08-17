package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dragonrider23/inca/common"
)

var Conn *sql.DB
var Ready = false

func Prepare() error {
	c := common.Config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Address,
		c.Database.Port,
		c.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		errors.New("Failed to open database")
	}
	Conn = db
	Ready = true
	return nil
}

func Close() {
	Conn.Close()
	return
}
