package main

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sort"
	"strconv"

	"github.com/xmdhs/creditget/db"
	"github.com/xmdhs/creditget/db/mysql"
	"github.com/xmdhs/creditget/model"
	"golang.org/x/exp/maps"
)

type config struct {
	Points map[string]string `json:"points"`
	DBUrl  string            `json:"dBUrl"`
}

var fieldName = map[string]string{
	"uid":         "uid",
	"name":        "用户名",
	"credits":     "总积分",
	"oltime":      "在线时间",
	"posts":       "回帖数",
	"threads":     "主题数",
	"friends":     "好友数",
	"medal":       "勋章数",
	"lastview":    "上次访问时间",
	"sex":         "性别",
	"groupname":   "用户组",
	"extgroupids": "扩展用户组",
}

func readConfig() config {
	c := config{}
	b, err := os.ReadFile(`config.json`)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}
	return c
}

var needSort = []string{"credits", "oltime", "posts", "threads", "friends", "medal",
	"extcredits1", "extcredits2", "extcredits3", "extcredits4", "extcredits5", "extcredits6", "extcredits7", "extcredits8"}

func main() {
	c := readConfig()
	fmap := maps.Clone(fieldName)
	maps.Copy(fmap, c.Points)
	cxt := context.Background()

	mysql, err := mysql.NewMysql(cxt, c.DBUrl)
	if err != nil {
		panic(err)
	}

	GetGroupSum(cxt, mysql)

	p, table := printTableName(model.CreditInfo{}, fmap)

	for _, v := range needSort {
		output(cxt, fmap[v], v, table, mysql, p, true)
	}
	output(cxt, "分积总", "credits", table, mysql, p, true)

	f, err := os.Create("一些统计.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("有效账号/总爬取账号：" + strconv.Itoa(must(mysql.GetAvailableUserSum(cxt))) + "/" + strconv.Itoa(must(mysql.GetSum(cxt))) + "  \n")
	f.WriteString("\n以下数据均为去除无效账号后的  \n")

	f.WriteString("零分：" + strconv.Itoa(must(mysql.GetNilSum(cxt, "credits"))) + "  \n")
	f.WriteString("零发帖：" + strconv.Itoa(must(mysql.GetNilSum(cxt, "threads"))) + "  \n")
	f.WriteString("零回帖：" + strconv.Itoa(must(mysql.GetNilSum(cxt, "posts"))) + "  \n")
	f.WriteString("零在线时间：" + strconv.Itoa(must(mysql.GetNilSum(cxt, "oltime"))) + "  \n")
	f.WriteString("零好友：" + strconv.Itoa(must(mysql.GetNilSum(cxt, "friends"))) + "  \n")

}

func output(cxt context.Context, name, field, table string, db db.DB, p *creditPrint, desc bool) {
	f, err := os.Create(name + ".md")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	defer bw.Flush()
	bw.WriteString(table)

	i := 0
	for offset := 0; offset < 2000; offset += 500 {
		c, err := db.GetRanks(cxt, field, 500, offset, desc)
		if err != nil {
			panic(err)
		}
		for _, v := range c {
			i++
			s := p.creditInfo2string(v, change, i)
			bw.WriteString(s)
			bw.WriteByte('\n')
		}
	}
}

func must[K any](k K, err error) K {
	if err != nil {
		panic(err)
	}
	return k
}

func GetGroupSum(cxt context.Context, db db.DB) {
	nl, err := db.GetGroupname(cxt)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(`组人数统计（不精准）.md`)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	type g struct {
		g string
		i int
	}
	gl := []g{}
	for _, v := range nl {
		if v == "" {
			continue
		}
		i, err := db.GetGroupCount(cxt, v)
		if err != nil {
			panic(err)
		}
		gl = append(gl, g{
			g: v,
			i: i,
		})
	}
	sort.Slice(gl, func(i, j int) bool {
		return gl[i].i > gl[j].i
	})
	for _, v := range gl {
		f.WriteString(v.g + ": ")
		f.WriteString(strconv.Itoa(v.i))
		f.WriteString("  \n")
	}
}
