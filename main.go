package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/xmdhs/creditget/db"
	"github.com/xmdhs/creditget/model"
	"github.com/xmdhs/creditget/profile"
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
		go toget(cxt, i, &w, mysql)
		t++
		if t > thread {
			w.Wait()
			t = 0
			mysql.InsterConfig(cxt, &model.Confing{
				ID:    0,
				VALUE: strconv.Itoa(i),
			})
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		}
	}

}

var c = &http.Client{
	Timeout: 10 * time.Second,
}

func toget(cxt context.Context, uid int, wait *sync.WaitGroup, db *db.MysqlDb) {
	defer wait.Done()
	var p *model.CreditInfo
	err := retry.Do(func() error {
		var err error
		p, err = profile.GetCredit(uid, c)
		return err
	}, getRetryOpts(20)...)
	if err != nil {
		panic(err)
	}
	err = retry.Do(func() error {
		return db.InsterCreditInfo(cxt, p)
	}, getRetryOpts(20)...)
	if err != nil {
		panic(err)
	}
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

func getRetryOpts(attempts uint) []retry.Option {
	if attempts == 0 {
		attempts = 15
	}
	return []retry.Option{
		retry.Attempts(attempts),
		retry.Delay(time.Second * 3),
		retry.LastErrorOnly(true),
		retry.MaxDelay(5 * time.Minute),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("retry %d: %v", n, err)
		}),
	}
}
