package get

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/xmdhs/creditget/sql"
)

type Friend struct {
	Ch         chan struct{}
	Wg         sync.WaitGroup
	SleepTime  int
	ProfileAPI string
}

func NewFriend(thread, SleepTime int, ProfileAPI string) *Friend {
	return &Friend{
		Ch:         make(chan struct{}, thread),
		Wg:         sync.WaitGroup{},
		SleepTime:  SleepTime,
		ProfileAPI: ProfileAPI,
	}
}

func (f *Friend) Friend(i int, uid string) {
	defer func() {
		time.Sleep(time.Duration(f.SleepTime) * time.Millisecond)
		<-f.Ch
		f.Wg.Done()
	}()
	u, uu := Getinfo(uid, f.ProfileAPI)

	uidi, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		log.Panicln(err)
		return
	}
	sql.Saveuserinfo(u, int(uidi))
	sql.Store(uid, uu.Name, uu.Friendstring, i+1)
}

func (f *Friend) Add(layers int) {
	i := 0
	for {
		if i > layers {
			break
		}
		lists := sql.GetList(i)
		for _, v := range lists {
			if sql.Find(v) {
				if v == "" {
					continue
				}
				f.Ch <- struct{}{}
				f.Wg.Add(1)
				go f.Friend(i, v)
			}
		}
		f.Wg.Wait()
		i++
	}
}
