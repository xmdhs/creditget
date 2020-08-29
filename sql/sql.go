package sql

import (
	"log"

	"github.com/xmdhs/creditget/get"
)

func Sqlget(id int) int {
	stmt, err := db.Prepare(`SELECT i FROM config WHERE id = ?`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query(id)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	rows.Next()
	var fid int
	rows.Scan(&fid)
	return fid
}

func Sqlup(s, id int) {
	stmt, err := db.Prepare("UPDATE config SET i = ? WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	stmt.Exec(s, id)
}

func Sqlinsert(id, start int) {
	_, err := db.Exec("INSERT INTO config VALUES (?,?)", id, start)
	if err != nil {
		log.Println(err)
	}

}

func Saveuserinfo(u get.Userinfo, uid int) {
	_, err := db.Exec(`INSERT INTO MCBBS VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		uid,
		u.Variables.Space.Username,
		u.Variables.Space.Credits,
		u.Variables.Space.Extcredits1,
		u.Variables.Space.Extcredits2,
		u.Variables.Space.Extcredits3,
		u.Variables.Space.Extcredits4,
		u.Variables.Space.Extcredits5,
		u.Variables.Space.Extcredits6,
		u.Variables.Space.Extcredits7,
		u.Variables.Space.Extcredits8,
		u.Variables.Space.Oltime,
		u.Variables.Space.Groupid,
		u.Variables.Space.Posts,
		u.Variables.Space.Threads,
		u.Variables.Space.Friends,
		u.Variables.Space.Views,
		u.Variables.Space.Adminid,
		getmedals(u.Variables.Space.Medals),
		u.Variables.Space.Digestposts,
		u.Variables.Space.Blacklist,
		u.Variables.Space.Emailstatus,
		u.Variables.Space.Lastvisit,
		u.Variables.Space.Avatarstatus,
		u.Variables.Space.Allowadmincp,
		u.Variables.Space.Extgroupids,
	)
	if err != nil {
		panic(err)
	}
}

func getmedals(medals interface{}) int {
	switch medals := medals.(type) {
	case []interface{}:
		return len(medals)
	case map[string]interface{}:
		return len(medals)
	default:
		return 0
	}
}
