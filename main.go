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
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
	w.Done()
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
