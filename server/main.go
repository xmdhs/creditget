package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/creditget/db"
	"github.com/xmdhs/creditget/db/cache"
	"github.com/xmdhs/creditget/db/mysql"
	"github.com/xmdhs/creditget/model"
	"github.com/xmdhs/creditget/profile"
)

func main() {
	mysqlDsn := os.Getenv("DSN")
	port := os.Getenv("PORT")

	cxt := context.Background()
	db, err := mysql.NewMysql(cxt, mysqlDsn)
	if err != nil {
		panic(err)
	}
	cacheDB := cache.NewMemCache(50000000, db)

	mux := httprouter.New()
	mux.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
			header.Set("Access-Control-Allow-Origin", "*")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	mux.GET("/userinfo/:uid", UserInfo(cacheDB))
	mux.GET("/rank/:uid/:field", rankHandler(cacheDB))

	s := http.Server{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      60 * time.Second,
		Addr:              ":" + port,
		Handler:           cors(mux),
	}
	s.ListenAndServe()
}

type ApiRep[V any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data V      `json:"data"`
}

func UserInfo(db db.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		s := r.FormValue("now")
		uid := p.ByName("uid")
		cxt := r.Context()
		uidi, err := strconv.Atoi(uid)
		if err != nil {
			handleErr(w, model.ApiErrInput, 400, err)
			return
		}
		var c *model.CreditInfo
		if s == "true" {
			c, err = profile.GetCredit(cxt, uidi, &http.Client{Timeout: 10 * time.Second})
		} else {
			c, err = db.GetCreditInfo(cxt, uidi)
		}
		if err != nil {
			handleErr(w, model.ApiDateBaseFail, 500, err)
			return
		}
		b, _ := json.Marshal(model.ApiRep[model.CreditInfo]{Data: *c})
		w.Header().Set("Cache-Control", "max-age=3600")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(b)
	}
}

func rankHandler(db db.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		uid := p.ByName("uid")
		uidi, err := strconv.Atoi(uid)
		if err != nil {
			handleErr(w, model.ApiErrInput, 400, err)
			return
		}
		cxt := r.Context()
		field := p.ByName("field")
		if field == "" {
			handleErr(w, model.ApiErrInput, 400, errors.New("field 不得为空"))
			return
		}
		rank, err := db.GetRank(cxt, uidi, field)
		if err != nil {
			handleErr(w, model.ApiDateBaseFail, 500, err)
			return
		}
		b, _ := json.Marshal(model.ApiRep[int]{Data: rank})
		w.Header().Set("Cache-Control", "max-age=3600")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(b)
	}
}

func handleErr(w http.ResponseWriter, code model.ApiErr, httpCode int, err error) {
	e := model.ApiRep[any]{}
	e.Code = int(code)
	e.Msg = err.Error()
	b, _ := json.Marshal(e)
	http.Error(w, string(b), httpCode)
	log.Println(err)
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Private-Network", "true")
		h.ServeHTTP(w, r)
	})
}
