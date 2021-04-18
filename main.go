package main

import (
	"encoding/json"
	"io/ioutil"
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
)

func main() {
	if len(os.Args) != 1 {
		output.GenAll()
	} else {
		var w sync.WaitGroup
		i := sql.Sqlget(0)
		if i < start {
			i = start
		}
		if i == 0 {
			i = 1
			sql.Sqlinsert(0, 1)
		}
		t := 0
		for ; i < end; i++ {
			w.Add(1)
			go toget(i, &w)
			t++
			if t > thread {
				w.Wait()
				t = 0
				sql.Sqlup(0, i)
			}
		}
	}
}

func toget(uid int, wait *sync.WaitGroup) {
	u := get.Getinfo(strconv.Itoa(uid))
	sql.Saveuserinfo(u, uid)
	log.Println(u.Variables.Space.Username, uid, u.Variables.Space.Credits)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	wait.Done()
}

func init() {
	readConfig()
}

func readConfig() {
	config := make(map[string]interface{}, 0)
	f, err := os.Open(`config.json`)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}
	start = int(config["start"].(float64))
	end = int(config["end"].(float64))
	thread = int(config["thread"].(float64))
	sleepTime = int(config["sleepTime"].(float64))
	get.ProfileAPI = config["disucuzApiAddress"].(string)
	for _, k := range output.Extcredits {
		if v, ok := config[k]; ok {
			output.Gendata[k] = v.(string)
		}
	}
}
