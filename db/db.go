package db

import (
	"database/sql"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("mysql", "sql:password@tcp(127.0.0.1:3306)/bookmanager")
}
