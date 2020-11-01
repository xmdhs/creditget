package main

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/xmdhs/creditget/get"
	"github.com/xmdhs/creditget/output"
	"github.com/xmdhs/creditget/sql"
)

const (
	start = 0
	// api/mobile/index.php?version=4&module=check 可获取论坛总人数
	end    = 3600000
	thread = 8
)

var w sync.WaitGroup

func main() {
	if len(os.Args) != 1 {
		output.GenAll()
	} else {
		a := end / thread
		for i := 0; i < thread; i++ {
			w.Add(1)
			go toget(start+a*i, start+a*(i+1), i)
		}
		w.Wait()
	}
}

func toget(s, end, id int) {
	a := sql.Sqlget(id)
	if a == 0 {
		sql.Sqlinsert(id, a)
		a = s
	}
	for i := start + a + 1; i <= start+end; i++ {
		u := get.Getinfo(strconv.Itoa(i))
		sql.Saveuserinfo(u, i)
		sql.Sqlup(i, id)
		log.Println(u.Variables.Space.Username, i, u.Variables.Space.Credits)
		time.Sleep(500 * time.Millisecond)
	}
	w.Done()
}
