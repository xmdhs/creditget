package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"

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
	sw.WriteString("| ")
	mt := reflect.TypeOf(m)
	fl := mt.NumField()
	for i := 0; i < fl; i++ {
		f := mt.Field(i)
		k := f.Tag.Get("json")
		v := fmap[k]
		if v == "" {
			v = k
		}
		c.k[i] = v
		sw.WriteString(v)
		if i < fl-1 {
			sw.WriteString(" | ")
		}
	}
	sw.WriteString(" |")
	return &c, sw.String()
}

func (c *creditPrint) creditInfo2string(m model.CreditInfo, change func(v any, field string) string) string {
	mv := reflect.ValueOf(m)
	fl := mv.NumField()
	sw := strings.Builder{}
	sw.WriteString("| ")
	for i := 0; i < fl; i++ {
		v := mv.Field(i)
		f := c.k[i]
		av := v.Interface()
		rv := change(av, f)
		if rv != "" {
			continue
		}
		sw.WriteString(escape(fmt.Sprint(av)))
		sw.WriteString(" | ")
	}
	sw.WriteString(" |")
	return sw.String()
}

var cnLoc, _ = time.LoadLocation("Asia/Shanghai")

func Lastview(v any, f string) string {
	if f != "lastview" {
		return ""
	}
	i := v.(int64)
	t := time.Unix(i, 0)
	t = t.In(cnLoc)
	return t.Format("2006-01-02 15:04:05")
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
