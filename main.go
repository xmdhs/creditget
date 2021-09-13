package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/xmdhs/creditget/get"
	"github.com/xmdhs/creditget/output"
	"github.com/xmdhs/creditget/sql"
)

var (
	// api/mobile/index.php?version=4&module=check 可获取论坛总人数
	start     int
	end       int
	thread    int
	sleepTime int = 500

	fast       bool
	fastUid    int = 1
	fastlayers int = 7

	profileAPI string
)

func main() {
	if len(os.Args) != 1 {
		output.GenAll()
	} else {
		var w sync.WaitGroup
		i := sql.Sqlget(0)
		if i == 0 {
			i = 1
			sql.Sqlinsert(0, 1)
		}
		if i < start {
			i = start
		}
		t := 0
		if !fast {
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
		} else {
			f := get.NewFriend(thread, sleepTime, profileAPI)
			f.Wg.Add(1)
			f.Ch <- struct{}{}
			f.Friend(-1, strconv.Itoa(fastUid))
			f.Add(fastlayers)
		}
	}
}

func toget(uid int, wait *sync.WaitGroup, profileAPI string) {
	u, _ := get.Getinfo(strconv.Itoa(uid), profileAPI)
	sql.Saveuserinfo(u, uid)
	log.Println(u.Variables.Space.Username, uid, u.Variables.Space.Credits)
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
	fast = c.Fast.On
	fastUid = c.Fast.UID
	fastlayers = c.Fast.Layers
	profileAPI = c.DisucuzAPIAddress
	for _, v := range output.GetExtcredits() {
		if pname, ok := c.Points[v]; ok {
			output.Gendata[v] = pname
		}
	}
}

type config struct {
	DisucuzAPIAddress string            `json:"disucuzApiAddress"`
	End               int               `json:"end"`
	Points            map[string]string `json:"points"`
	Fast              configFast        `json:"fast"`
	SleepTime         int               `json:"sleepTime"`
	Start             int               `json:"start"`
	Thread            int               `json:"thread"`
}

type configFast struct {
	Layers int  `json:"layers"`
	On     bool `json:"on"`
	UID    int  `json:"uid"`
}
