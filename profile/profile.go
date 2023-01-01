package profile

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/xmdhs/creditget/model"
)

type ErrHttpCode int

func (e ErrHttpCode) Error() string {
	return fmt.Sprintf("http code: %v", int(e))
}

func GetCredit(uid int, c *http.Client) (*model.CreditInfo, error) {
	req, err := http.NewRequest("GET", "https://www.mcbbs.net/home.php?mod=space&uid="+strconv.Itoa(uid), nil)
	if err != nil {
		return nil, fmt.Errorf("GetCredit: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	reps, err := c.Do(req)
	if reps != nil {
		defer reps.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("GetCredit: %w", err)
	}
	if reps.StatusCode != 200 {
		return nil, fmt.Errorf("GetCredit: %w", ErrHttpCode(reps.StatusCode))
	}

	d, err := goquery.NewDocumentFromReader(reps.Body)
	if err != nil {
		return nil, fmt.Errorf("GetCredit: %w", err)
	}
	m := model.CreditInfo{}

	m.Name = getName(d)
	m.Uid = int32(uid)

	cl := getCredits(d)
	m.Credits = cl[0]
	m.Extcredits1 = cl[1]
	m.Extcredits2 = cl[2]
	m.Extcredits3 = cl[3]
	m.Extcredits4 = cl[4]
	m.Extcredits5 = cl[5]
	m.Extcredits6 = cl[6]
	m.Extcredits7 = cl[7]
	m.Extcredits8 = cl[8]

	m.Friends = getFriends(d)
	m.Posts = getPosts(d)
	m.Threads = getThreads(d)
	m.Oltime = getOltime(d)
	m.Groupname = getGroupname(d)
	m.Medal = getMedal(d)
	m.Lastview = getLastview(d)
	m.Extgroupids = getExtgroupids(d)
	m.Sex = getSex(d)

	return &m, nil
}
