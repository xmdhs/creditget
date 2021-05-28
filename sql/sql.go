package sql

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-sqlite3"
	j "github.com/xmdhs/creditget/json"
)

func Sqlget(id int) int {
	stmt, err := db.Prepare(`SELECT i FROM config WHERE id = ?`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	rows.Next()
	var fid int
	rows.Scan(&fid)
	return fid
}

func Sqlup(id, s int) {
	stmt, err := db.Prepare("UPDATE config SET i = ? WHERE id = ?")
	if err != nil {
		e := sqlite3.Error{}
		if errors.As(err, &e) {
			if e.Code == sqlite3.ErrBusy || e.Code == sqlite3.ErrLocked {
				log.Println(err)
				time.Sleep(1 * time.Second)
				Sqlup(id, s)
				return
			}
		}
		panic(err)
	}
	defer stmt.Close()
	stmt.Exec(s, id)
}

func Sqlinsert(id, start int) {
	_, err := db.Exec("INSERT INTO config VALUES (?,?)", id, start)
	if err != nil {
		e := sqlite3.Error{}
		if errors.As(err, &e) {
			if e.Code == sqlite3.ErrConstraint {
				log.Println(err)
				return
			}
			if e.Code == sqlite3.ErrBusy || e.Code == sqlite3.ErrLocked {
				log.Println(err)
				time.Sleep(1 * time.Second)
				Sqlinsert(id, start)
				return
			}
		}
		panic(err)
	}
}

func Saveuserinfo(u j.Userinfo, uid int) {
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
		u.Variables.Space.Group.Grouptitle,
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
		e := sqlite3.Error{}
		if errors.As(err, &e) {
			if e.Code == sqlite3.ErrConstraint {
				log.Println(err)
				return
			}
			if e.Code == sqlite3.ErrBusy || e.Code == sqlite3.ErrLocked {
				log.Println(err)
				time.Sleep(1 * time.Second)
				Saveuserinfo(u, uid)
				return
			}
		}
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

func Find(uid string) bool {
	rows, err := db.Query("SELECT uid FROM friend WHERE uid=?", uid)
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		return false
	}
	return true
}

var one *sync.Once = &sync.Once{}
var stmt *sql.Stmt

func Store(uid, name, friend string, i int) {
	one.Do(func() {
		var err error
		stmt, err = db.Prepare("INSERT INTO friend VALUES (?,?,?,?)")
		if err != nil {
			log.Println(err)
			return
		}
	})
	_, err := stmt.Exec(uid, name, friend, i)
	if err != nil {
		log.Println(err)
		return
	}
}

func GetList(i int) []string {
	lists := make([]string, 0)
	rows, err := db.Query("SELECT * FROM friend WHERE i=?", i)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var uidd string
	var name string
	var friend string
	for rows.Next() {
		err := rows.Scan(&uidd, &name, &friend, &i)
		if err != nil {
			panic(err)
		}
		list := strings.Split(friend, ",")
		lists = append(lists, list...)
	}
	return lists
}
