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

	fast       bool
	fastUid    int = 1
	fastlayers int = 7
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
				go toget(i, &w)
				t++
				if t > thread {
					w.Wait()
					t = 0
					sql.Sqlup(0, i+1)
				}
			}
		} else {
			get.Wg.Add(1)
			get.Ch <- struct{}{}
			get.Friend(-1, strconv.Itoa(fastUid))
			get.Add(fastlayers)
		}
	}
}

func toget(uid int, wait *sync.WaitGroup) {
	u, _ := get.Getinfo(strconv.Itoa(uid))
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
	if err != nil {
		panic(err)
	}
	defer f.Close()
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
	fast = config["fast"].(bool)
	fastUid = int(config["fast_uid"].(float64))
	fastlayers = int(config["fast_layers"].(float64))
	for _, k := range output.Extcredits {
		if v, ok := config[k]; ok {
			output.Gendata[k] = v.(string)
		}
	}
}
