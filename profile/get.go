package profile

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "time/tzdata"

	"github.com/PuerkitoBio/goquery"
)

func getName(d *goquery.Document) string {
	return strings.TrimSpace(d.Find("#uhd > div.h.cl > h2").Text())
}

var numReg = regexp.MustCompile(`-?\d+`)

func getFriends(d *goquery.Document) int32 {
	f := d.Find("#ct li > a:contains(好友数)").Text()
	return toInt32(numReg.FindString(f))
}

func getPosts(d *goquery.Document) int32 {
	f := d.Find("#ct li > a:contains(回帖数)").Text()
	return toInt32(numReg.FindString(f))
}

func getThreads(d *goquery.Document) int32 {
	f := d.Find("#ct li > a:contains(主题数)").Text()
	return toInt32(numReg.FindString(f))
}

func toInt32(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int32(i)
}

func getCredits(d *goquery.Document) [9]int32 {
	cl := [9]int32{}
	has := false
	f := d.Find("#psts > ul > li")
	i := 0
	f.EachWithBreak(func(_ int, s *goquery.Selection) bool {
		t := s.Text()
		if strings.Contains(t, "积分") {
			has = true
		}
		if !has {
			return true
		}
		n := numReg.FindString(t)
		cl[i] = toInt32(n)
		i++
		return true
	})
	return cl
}

func getOltime(d *goquery.Document) int32 {
	f := d.Find("#pbbs > li:has(em:contains(在线时间))")
	return toInt32(numReg.FindString(f.Text()))
}

func getGroupname(d *goquery.Document) string {
	f := d.Find("ul > li:has(em.xg1:contains(用户组):not(:contains(扩展用户组))) > span > a")
	return f.Text()
}

func getMedal(d *goquery.Document) int32 {
	s := d.Find(".md_ctrl > a[href='home.php?mod=medal'] > img")
	return int32(s.Length())
}

var timeReg = regexp.MustCompile(`\d{4}-\d{1,2}-\d{1,2} \d{2}:\d{2}`)
var shanhai, _ = time.LoadLocation("Asia/Shanghai")

func getLastview(d *goquery.Document) int64 {
	f := d.Find("#pbbs > li:has(em:contains(最后访问))")
	ts := timeReg.FindString(f.Text())
	t, err := time.ParseInLocation("2006-1-2 15:04", ts, shanhai)
	if err != nil {
		return 0
	}
	return t.Unix()
}

func getExtgroupids(d *goquery.Document) string {
	f := d.Find("ul > li:has(em.xg1:contains(扩展用户组))")
	f.Find("em").Remove()
	return f.Text()
}

func getSex(d *goquery.Document) int32 {
	f := d.Find("#ct ul.pf_l.cl	> li:has(em:contains(性别))")
	f.Find("em").Remove()
	find := f.Text()

	switch {
	case strings.Contains(find, "男"):
		return 1
	case strings.Contains(find, "女"):
		return 2
	default:
		return 0
	}

}
