package get

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var c = http.Client{Timeout: 5 * time.Second}

func h(uid string) ([]byte, error) {
	reqs, err := http.NewRequest("GET", `https://www.mcbbs.net/api/mobile/index.php?version=4&module=profile&uid=`+uid, nil)
	if err != nil {
		return nil, err
	}
	reqs.Header.Set("Accept", "*/*")
	reqs.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	rep, err := c.Do(reqs)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if rep.StatusCode != http.StatusOK {
		return nil, errors.New(rep.Status)
	}
	b, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func json2userinfo(b []byte) (Userinfo, error) {
	u := Userinfo{}
	err := json.Unmarshal(b, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func Getinfo(uid string) Userinfo {
	for {
		b, err := h(uid)
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		u, err := json2userinfo(b)
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		return u
	}
}
