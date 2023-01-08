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

	"github.com/xmdhs/creditget/db"
	"github.com/xmdhs/creditget/db/mysql"
	"github.com/xmdhs/creditget/model"
)

func main() {
	mysqlDsn := os.Getenv("DSN")
	port := os.Getenv("PORT")

	cxt := context.Background()
	db, err := mysql.NewMysql(cxt, mysqlDsn)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/userinfo", UserInfo(db))
	mux.HandleFunc("/rank", rankHandler(db))

	s := http.Server{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      60 * time.Second,
		Addr:              ":" + port,
		Handler:           mux,
	}
	s.ListenAndServe()
}

type ApiRep[V any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data V      `json:"data"`
}

func UserInfo(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.FormValue("uid")
		cxt := r.Context()
		uidi, err := strconv.Atoi(uid)
		if err != nil {
			handleErr(w, model.ApiErrInput, 400, err)
			return
		}
		c, err := db.GetCreditInfo(cxt, uidi)
		if err != nil {
			handleErr(w, model.ApiDateBaseFail, 500, err)
			return
		}
		b, _ := json.Marshal(model.ApiRep[model.CreditInfo]{Data: *c})
		w.Header().Set("Cache-Control", "max-age=3600")
		w.Write(b)
	}
}

func rankHandler(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.FormValue("uid")
		uidi, err := strconv.Atoi(uid)
		if err != nil {
			handleErr(w, model.ApiErrInput, 400, err)
			return
		}
		cxt := r.Context()
		field := r.FormValue("field")
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
