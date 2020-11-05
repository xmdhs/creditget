package output

import (
	"bufio"
	"database/sql"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Sqlget(key string, limit int, desc bool) string {
	sw := strings.Builder{}
	var word string
	if desc {
		word = "DESC"
	} else {
		word = "ASC"
	}
	stmt, err := db.Prepare(`SELECT * FROM MCBBS ORDER BY ` + key + ` ` + word + ` LIMIT ?`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query(limit)
	defer rows.Close()
	text, i := setTableName()
	sw.WriteString(text)
	sw.WriteString(genSeparate(i))
	u := rowsget(rows)
	for i, v := range u {
		i++
		sw.WriteString("| " + strconv.Itoa(i) + " | ")
		sw.WriteString(v.UID + " | ")
		sw.WriteString(v.Username + " | ")
		sw.WriteString(v.Credits + " | ")
		sw.WriteString(genTableValue("extcredits1", v.Extcredits1))
		sw.WriteString(genTableValue("extcredits2", v.Extcredits2))
		sw.WriteString(genTableValue("extcredits3", v.Extcredits3))
		sw.WriteString(genTableValue("extcredits4", v.Extcredits4))
		sw.WriteString(genTableValue("extcredits5", v.Extcredits5))
		sw.WriteString(genTableValue("extcredits6", v.Extcredits6))
		sw.WriteString(genTableValue("extcredits7", v.Extcredits7))
		sw.WriteString(genTableValue("extcredits8", v.Extcredits8))
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

var Gendata = map[string]string{
	"credits":     "总积分",
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
	wait := sync.WaitGroup{}
	for k, v := range Gendata {
		wait.Add(1)
		k, v := k, v
		go func() {
			s := Sqlget(k, 1000, true)
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
			wait.Done()
		}()
	}
	wait.Wait()
	f1, err := os.Create(`分积总.md`)
	defer f1.Close()
	if err != nil {
		panic(err)
	}
	s := Sqlget("credits", 1000, false)
	_, err = f1.WriteString(s)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(`组人数统计（不精准）.txt`)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	m := GetGroupSum()
	bf := bufio.NewWriter(f)
	defer bf.Flush()
	for k, v := range m {
		_, err := bf.WriteString(k + ": " + strconv.Itoa(v) + "\n")
		if err != nil {
			panic(err)
		}
	}
	f3, err := os.Create(`一些统计.txt`)
	defer f3.Close()
	f2 := bufio.NewWriter(f3)
	defer f2.Flush()
	f2.WriteString("有效账号/总爬取账号：" + strconv.Itoa(GetAvailableUserSum()) + "/" + strconv.Itoa(GetSum()) + "\n")
	f2.WriteString("\n以下数据均为去除无效账号后的\n")
	f2.WriteString("未设置邮箱：" + strconv.Itoa(GetNotEmailsSum()) + "\n")
	f2.WriteString("未设置头像：" + strconv.Itoa(GetNotSetAvatarSum()) + "\n")
	f2.WriteString("零分：" + strconv.Itoa(GetNilCreditsSum()) + "\n")
	f2.WriteString("零发帖：" + strconv.Itoa(GetNilPosts()) + "\n")
	f2.WriteString("零回帖：" + strconv.Itoa(GetNilThreads()) + "\n")
	f2.WriteString("零在线时间：" + strconv.Itoa(GetNilOltime()) + "\n")
}

func GetGroupSum() map[string]int {
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
		rows := db.QueryRow("SELECT COUNT(*) FROM mcbbs WHERE groupname ='" + v + "' OR extgroupids LIKE '%" + v + "%';")
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

func GetNotSetAvatarSum() int {
	return getNilSum(`Avatarstatus`)
}

func GetNotEmailsSum() int {
	return getNilSum(`emailstatus`)
}

func GetNilCreditsSum() int {
	return getNilSum(`credits`)
}

func GetNilPosts() int {
	return getNilSum(`posts`)
}

func GetNilThreads() int {
	return getNilSum(`threads`)
}

func GetNilOltime() int {
	return getNilSum(`oltime`)
}

func getNilSum(name string) int {
	rows := db.QueryRow(`SELECT COUNT(*) FROM mcbbs WHERE ` + name + ` = 0 AND NOT lastactivitydb = 0`)
	i := 0
	err := rows.Scan(&i)
	if err != nil {
		panic(err)
	}
	return i
}

func GetSum() int {
	rows := db.QueryRow(`SELECT COUNT(*) FROM mcbbs`)
	i := 0
	err := rows.Scan(&i)
	if err != nil {
		panic(err)
	}
	return i
}

func GetAvailableUserSum() int {
	rows := db.QueryRow(`SELECT COUNT(*) FROM mcbbs WHERE NOT lastactivitydb = 0`)
	i := 0
	err := rows.Scan(&i)
	if err != nil {
		panic(err)
	}
	return i
}
