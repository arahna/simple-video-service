package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase(dataSourceName string) (*sql.DB, error) {
	return sql.Open("mysql", dataSourceName)
}
