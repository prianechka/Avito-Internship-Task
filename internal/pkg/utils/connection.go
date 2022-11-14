package utils

import (
	"Avito-Internship-Task/configs"
	"database/sql"
	"fmt"
)

func NewMySQLConnction(conn configs.MySQLConnectionParams) *sql.DB {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conn.User, conn.Password,
		conn.Host, conn.Port, conn.Database)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err)
	}
	return db
}
