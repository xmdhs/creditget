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
	u := Userinfo{
		Variables: variables{
			Space: space{
				Adminid:      "0",
				Allowadmincp: "0",
				Avatarstatus: "0",
				Blacklist:    "0",
				Credits:      "0",
				Digestposts:  "0",
				Emailstatus:  "0",
				Extcredits1:  "0",
				Extcredits2:  "0",
				Extcredits3:  "0",
				Extcredits4:  "0",
				Extcredits5:  "0",
				Extcredits6:  "0",
				Extcredits7:  "0",
				Extcredits8:  "0",
				Extgroupids:  "0",
				Friends:      "0",
				Groupid:      "0",
				Lastvisit:    "0",
				Oltime:       "0",
				Posts:        "0",
				Threads:      "0",
				UID:          "0",
				Username:     "0",
				Views:        "0",
				Medals:       nil,
			},
		},
	}
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
