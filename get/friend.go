package get

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xmdhs/creditget/sql"
)

var Ch = make(chan struct{}, 10)
var Wg = &sync.WaitGroup{}

func Friend(i int, uid string) {
	defer func() {
		time.Sleep(500 * time.Millisecond)
		<-Ch
		Wg.Done()
	}()
	if i > 7 {
		return
	}
	u, uu := Getinfo(uid)

	uidi, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		panic(err)
	}
	sql.Saveuserinfo(u, int(uidi))
	sql.Store(uid, uu.Name, uu.Friendstring, i+1)
}

func Add(layers int) {
	i := 0
	for {
		if i > layers {
			break
		}
		lists := make([]string, 0)
		rows, err := sql.Getdb().Query("SELECT * FROM MCBBS WHERE i=?", i)
		if err != nil {
			rows.Close()
			panic(err)
		}
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
		rows.Close()
		for _, v := range lists {
			if sql.Find(v) {
				Ch <- struct{}{}
				Wg.Add(1)
				go Friend(i, v)
			}
		}
		Wg.Wait()
		i++
	}
}
