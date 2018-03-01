package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const dataSourceName = `root:1234@/video`

func InitDatabase() (*sql.DB, error) {
	return sql.Open("mysql", dataSourceName)
}
