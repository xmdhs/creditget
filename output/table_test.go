package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/xmdhs/creditget/model"
	"github.com/xmdhs/creditget/profile"
)

func Test_printTableName(t *testing.T) {
	_, s := printTableName(model.CreditInfo{}, fieldName)
	fmt.Println(s)
}

func Test_creditPrint_creditInfo2string(t *testing.T) {
	p, err := profile.GetCredit(context.Background(), 1770442, &http.Client{Timeout: 10 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	print, _ := printTableName(model.CreditInfo{}, fieldName)
	s := print.creditInfo2string(*p, change, 1)
	fmt.Println(s)
}
