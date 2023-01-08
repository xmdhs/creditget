package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/xmdhs/creditget/model"
)

type creditPrint struct {
	k map[int]string
}

func printTableName(m model.CreditInfo, fmap map[string]string) (*creditPrint, string) {
	c := creditPrint{
		k: map[int]string{},
	}
	sw := strings.Builder{}
	sw.WriteString("| 排名 |")
	mt := reflect.TypeOf(m)
	fl := mt.NumField()
	for i := 0; i < fl; i++ {
		f := mt.Field(i)
		k := f.Tag.Get("json")
		if k == "" {
			k = f.Name
		}
		c.k[i] = k
		v := fmap[k]
		if v == "" {
			v = k
		}
		sw.WriteString(v)
		if i < fl-1 {
			sw.WriteString(" | ")
		}
	}
	sw.WriteString(" |\n")
	sw.WriteString("|")
	for i := 0; i < fl+1; i++ {
		sw.WriteString(" - |")
	}
	sw.WriteByte('\n')
	return &c, sw.String()
}

func (c *creditPrint) creditInfo2string(m model.CreditInfo, change func(v any, field string) string, i int) string {
	mv := reflect.ValueOf(m)
	fl := mv.NumField()
	sw := strings.Builder{}
	sw.WriteString("| ")
	sw.WriteString(strconv.Itoa(i))
	sw.WriteString(" | ")
	for i := 0; i < fl; i++ {
		v := mv.Field(i)
		f := c.k[i]
		av := v.Interface()
		rv := change(av, f)
		if rv != "" {
			av = rv
		}
		sw.WriteString(escape(fmt.Sprint(av)))
		sw.WriteString(" | ")
	}
	sw.WriteString(" |")
	return sw.String()
}

var cnLoc, _ = time.LoadLocation("Asia/Shanghai")

func change(v any, f string) string {
	switch f {
	case "lastview":
		i := v.(int64)
		t := time.Unix(i, 0)
		t = t.In(cnLoc)
		return t.Format("2006-01-02 15:04:05")
	case "sex":
		i := v.(int32)
		if i == 1 {
			return "男"
		}
		if i == 2 {
			return "女"
		}
		return "保密"
	}
	return ""
}

var escapeWord = []string{
	`\`,
	"`",
	"*",
	"_",
	"{",
	"}",
	"[",
	"]",
	"(",
	")",
	"#",
	"+",
	"-",
	".",
	"!",
	"|",
}

func escape(w string) string {
	for _, v := range escapeWord {
		if strings.Contains(w, v) {
			w = strings.ReplaceAll(w, v, `\`+v)
		}
	}
	return w
}
