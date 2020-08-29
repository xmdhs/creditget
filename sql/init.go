package sql

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
	_, err = db.Exec(
		"CREATE TABLE 		IF NOT EXISTS MCBBS " +
			"(UID			INT	PRIMARY	KEY	NOT	NULL," +
			"NAME			TEX NOT NULL," +
			"credits		INT	NOT NULL," +
			"extcredits1 	INT	NOT	NULL," +
			"extcredits2 	INT	NOT	NULL," +
			"extcredits3 	INT	NOT	NULL," +
			"extcredits4 	INT	NOT	NULL," +
			"extcredits5 	INT	NOT	NULL," +
			"extcredits6 	INT	NOT	NULL," +
			"extcredits7 	INT	NOT	NULL," +
			"extcredits8 	INT	NOT	NULL," +
			"oltime			INT	NOT	NULL," +
			"groupid		INT	NOT	NULL," +
			"posts			INT	NOT	NULL," +
			"threads		INT	NOT	NULL," +
			"friends		INT	NOT	NULL," +
			"views			INT	NOT	NULL," +
			"adminid		INT	NOT	NULL," +
			"medal			INT	NOT	NULL," +
			"digestposts	INT	NOT	NULL," +
			"blacklist		INT	NOT	NULL," +
			"emailstatus	INT	NOT	NULL," +
			"lastvisit		TEXT	NOT	NULL," +
			"Avatarstatus	INT	NOT	NULL," +
			"Allowadmincp	INT	NOT	NULL," +
			"extgroupids 	TEXT NOT	NULL)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS config (id INT PRIMARY KEY NOT NULL,i INT NOT NULL)`)
	if err != nil {
		panic(err)
	}
}
