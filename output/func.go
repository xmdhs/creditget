package output

import (
	"bufio"
	"database/sql"
	"os"
	"strconv"
	"strings"
)

func Sqlget(key string, limit int) string {
	sw := strings.Builder{}
	stmt, err := db.Prepare(`SELECT * FROM MCBBS ORDER BY ` + key + ` DESC LIMIT ?`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query(limit)
	defer rows.Close()
	sw.WriteString("| 排名 | uid | 用户名 | 积分 | 人气 | 金粒 | 金锭 | 绿宝石 | 下界之星 | 贡献 | 爱心 | 钻石 | 勋章数 | 精华数 | 设置黑名单数 | 在线时间 | 回帖数 | 主题数 | 好友数 | 空间查看次数 | 用户组 | 扩展用户组 | 上次访问时间 |\n")
	sw.WriteString("| - | - | - | - | - | - | - | - | - | - | - | - | - | - | - | - | - | - | - | - | - | - | - |\n")
	u := rowsget(rows)
	for i, v := range u {
		i++
		sw.WriteString("| " + strconv.Itoa(i) + " | ")
		sw.WriteString(v.UID + " | ")
		sw.WriteString(v.Username + " | ")
		sw.WriteString(v.Credits + " | ")
		sw.WriteString(v.Extcredits1 + " | ")
		sw.WriteString(v.Extcredits2 + " | ")
		sw.WriteString(v.Extcredits3 + " | ")
		sw.WriteString(v.Extcredits4 + " | ")
		sw.WriteString(v.Extcredits5 + " | ")
		sw.WriteString(v.Extcredits6 + " | ")
		sw.WriteString(v.Extcredits7 + " | ")
		sw.WriteString(v.Extcredits8 + " | ")
		sw.WriteString(strconv.Itoa(v.Medals) + " | ")
		sw.WriteString(v.Digestposts + " | ")
		sw.WriteString(v.Blacklist + " | ")
		sw.WriteString(v.Oltime + " | ")
		sw.WriteString(v.Posts + " | ")
		sw.WriteString(v.Threads + " | ")
		sw.WriteString(v.Friends + " | ")
		sw.WriteString(v.Views + " | ")
		sw.WriteString(v.Group + " | ")
		sw.WriteString(v.Extgroupids + " | ")
		sw.WriteString(v.Lastvisit + " |\n")
	}
	return sw.String()
}

func rowsget(rows *sql.Rows) []userdata {
	ulist := make([]userdata, 0)
	for rows.Next() {
		u := userdata{}
		err := rows.Scan(&u.UID, &u.Username, &u.Credits, &u.Extcredits1, &u.Extcredits2, &u.Extcredits3, &u.Extcredits4, &u.Extcredits5, &u.Extcredits6, &u.Extcredits7, &u.Extcredits8,
			&u.Oltime, &u.Group, &u.Posts, &u.Threads, &u.Friends, &u.Views, &u.Adminid, &u.Medals, &u.Digestposts, &u.Blacklist, &u.Emailstatus, &u.Lastvisit, &u.Avatarstatus,
			&u.Allowadmincp, &u.Extgroupids)
		if err != nil {
			panic(err)
		}
		ulist = append(ulist, u)
	}
	return ulist
}

type userdata struct {
	Adminid      string
	Allowadmincp string
	Avatarstatus string
	Blacklist    string
	Credits      string
	Digestposts  string
	Emailstatus  string
	Extcredits1  string
	Extcredits2  string
	Extcredits3  string
	Extcredits4  string
	Extcredits5  string
	Extcredits6  string
	Extcredits7  string
	Extcredits8  string
	Extgroupids  string
	Friends      string
	Lastvisit    string
	Oltime       string
	Posts        string
	Threads      string
	UID          string
	Username     string
	Views        string
	Medals       int
	Group        string
}

var gendata = map[string]string{
	"credits":     "总积分",
	"extcredits1": "人气",
	"extcredits2": "金粒",
	"extcredits3": "金锭",
	"extcredits4": "绿宝石",
	"extcredits5": "下界之星",
	"extcredits6": "贡献",
	"extcredits7": "爱心",
	"extcredits8": "钻石",
	"oltime":      "在线时间",
	"posts":       "回帖数",
	"threads":     "主题数",
	"friends":     "好友数",
	"views":       "空间查看次数",
	"medal":       "勋章数",
	"digestposts": "精华数",
	"blacklist":   "设置的黑名单数",
}

func GenAll() {
	for k, v := range gendata {
		s := Sqlget(k, 100)
		f, err := os.Create(v + ".md")
		if err != nil {
			f.Close()
			panic(err)
		}
		_, err = f.WriteString(s)
		if err != nil {
			f.Close()
			panic(err)
		}
		f.Close()
	}
	f, err := os.Create(`组人数统计（不精准）.txt`)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	m := getGroupSum()
	bf := bufio.NewWriter(f)
	defer bf.Flush()
	for k, v := range m {
		_, err := bf.WriteString(k + ": " + strconv.Itoa(v) + "\n")
		if err != nil {
			panic(err)
		}
	}
}

func getGroupSum() map[string]int {
	rows, err := db.Query(`SELECT DISTINCT groupname FROM MCBBS`)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	list := make([]string, 0)
	for rows.Next() {
		var groupname string
		err := rows.Scan(&groupname)
		if err != nil {
			panic(err)
		}
		list = append(list, groupname)
	}
	m := make(map[string]int, 0)
	for _, v := range list {
		rows := db.QueryRow("SELECT COUNT(UID) FROM mcbbs WHERE groupname ='" + v + "' OR extgroupids LIKE '%" + v + "%';")
		if err != nil {
			panic(err)
		}
		var i int
		err := rows.Scan(&i)
		if err != nil {
			panic(err)
		}
		m[v] = i
	}
	return m
}
