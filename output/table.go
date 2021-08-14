package output

import (
	"strings"
)

func setTableName() (string, int) {
	sw := strings.Builder{}
	sw.WriteString("| 排名 | uid | 用户名 | 积分 | ")
	sum := 0
	for _, v := range extcredits {
		if v, ok := Gendata[v]; ok {
			sum++
			sw.WriteString(v)
			sw.WriteString(" | ")
		}
	}
	sw.WriteString("勋章数 | 精华数 | 设置黑名单数 | 在线时间 | 回帖数 | 主题数 | 好友数 | 空间查看次数 | 用户组 | 扩展用户组 | 上次访问时间 |\n")
	return sw.String(), sum
}

var extcredits = []string{
	"extcredits1",
	"extcredits2",
	"extcredits3",
	"extcredits4",
	"extcredits5",
	"extcredits6",
	"extcredits7",
	"extcredits8",
}

func GetExtcredits() []string {
	return extcredits
}

func genSeparate(i int) string {
	i = i + 15
	sw := strings.Builder{}
	for a := 0; a < i; a++ {
		sw.WriteString("| - ")
	}
	sw.WriteString("|\n")
	return sw.String()
}

func genTableValue(name, value string) string {
	if _, ok := Gendata[name]; ok {
		return value + " | "
	}
	return ""
}
