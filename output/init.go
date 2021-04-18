package output

import (
	"database/sql"

	//数据库驱动
	_ "github.com/mattn/go-sqlite3"
	asql "github.com/xmdhs/creditget/sql"
)

var db *sql.DB

func init() {
	db = asql.Getdb()
}
