package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/xmdhs/creditget/db"
)

var (
	start int
	// api/mobile/index.php?version=4&module=check 可获取论坛总人数
	end       int
	thread    int
	sleepTime int = 500

	DBUrl string
)

func main() {
	mysql, err := db.NewMysql(DBUrl)
	if err != nil {
		panic(err)
	}
	cxt := context.Background()

	var w sync.WaitGroup
	i := 1
	v, err := mysql.SelectConfig(cxt, 0)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	if v.VALUE != "" {
		i, err = strconv.Atoi(v.VALUE)
		if err != nil {
			panic(err)
		}
	}

	if i < start {
		i = start
	}

	t := 0
	for ; i < end; i++ {
		w.Add(1)
		go toget(i, &w, profileAPI)
		t++
		if t > thread {
			w.Wait()
			t = 0
			sql.Sqlup(0, i+1)
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		}
	}

}

var c = &http.Client{
	Timeout: 10 * time.Second,
}

func toget(uid int, wait *sync.WaitGroup) {

	wait.Done()
}

func init() {
	readConfig()
}

func readConfig() {
	c := config{}
	b, err := os.ReadFile(`config.json`)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}
	start = c.Start
	end = c.End
	thread = c.Thread
	sleepTime = c.SleepTime
	DBUrl = c.DBUrl
}

type config struct {
	End       int               `json:"end"`
	Points    map[string]string `json:"points"`
	SleepTime int               `json:"sleepTime"`
	Start     int               `json:"start"`
	Thread    int               `json:"thread"`
	DBUrl     string            `json:"dBUrl"`
}
