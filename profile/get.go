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

var numReg = regexp.MustCompile(`\d+`)

func getFriends(d *goquery.Document) int32 {
	f := d.Find("#ct ul.cl.bbda.pbm.mbm > li > a:nth-child(2)").Text()
	return toInt32(numReg.FindString(f))
}

func getPosts(d *goquery.Document) int32 {
	f := d.Find("#ct ul.cl.bbda.pbm.mbm > li > a:nth-child(4)").Text()
	return toInt32(numReg.FindString(f))
}

func getThreads(d *goquery.Document) int32 {
	f := d.Find("#ct ul.cl.bbda.pbm.mbm > li > a:nth-child(6)").Text()
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
	f := d.Find("#pbbs > li")
	oltime := ""
	f.EachWithBreak(func(i int, s *goquery.Selection) bool {
		t := s.Text()
		if strings.Contains(t, "在线时间") {
			oltime = t
			return false
		}
		return true
	})
	return toInt32(numReg.FindString(oltime))
}

func getGroupname(d *goquery.Document) string {
	ems := d.Find("em.xg1")
	findS := ems.First()

	ems.EachWithBreak(func(i int, s *goquery.Selection) bool {
		t := s.Text()
		if strings.Contains(t, "用户组") && !strings.Contains(t, "扩展用户组") {
			findS = s
			return false
		}
		return true
	})

	return findS.Parent().Find("a").Text()
}

func getMedal(d *goquery.Document) int32 {
	s := d.Find(".md_ctrl > a[href='home.php?mod=medal'] > img")
	return int32(s.Length())
}

var timeReg = regexp.MustCompile(`\d{4}-\d{1,2}-\d{1,2} \d{2}:\d{2}`)
var shanhai, _ = time.LoadLocation("Asia/Shanghai")

func getLastview(d *goquery.Document) int64 {
	f := d.Find("#pbbs > li")
	lastv := ""
	f.EachWithBreak(func(i int, s *goquery.Selection) bool {
		t := s.Text()
		if strings.Contains(t, "最后访问") {
			lastv = t
			return false
		}
		return true
	})
	ts := timeReg.FindString(lastv)
	t, err := time.ParseInLocation("2006-1-2 15:04", ts, shanhai)
	if err != nil {
		return 0
	}
	return t.Unix()
}

func getExtgroupids(d *goquery.Document) string {
	ems := d.Find("em.xg1")
	var findS *goquery.Selection

	ems.EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "扩展用户组") {
			findS = s
			return false
		}
		return true
	})
	if findS == nil {
		return ""
	}

	p := findS.Parent()
	findS.Remove()
	return p.Text()
}

func getSex(d *goquery.Document) int32 {
	f := d.Find("#ct > div > div.bm.bw0 > div > div.bm_c.u_profile > div:nth-child(1) > ul:nth-child(5) > li")
	find := f.First().Text()

	f.EachWithBreak(func(i int, s *goquery.Selection) bool {
		t := s.Text()
		if strings.Contains(t, "性别") {
			find = t
			return false
		}
		return true
	})

	switch {
	case strings.Contains(find, "男"):
		return 1
	case strings.Contains(find, "女"):
		return 2
	default:
		return 0
	}

}
