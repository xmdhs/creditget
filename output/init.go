package output

import (
	"database/sql"

	//数据库驱动
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./credit.db")
	if err != nil {
		panic(err)
	}
}
