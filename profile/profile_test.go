package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGetCredit(t *testing.T) {
	c := &http.Client{Timeout: 10 * time.Second}
	m, err := GetCredit(context.Background(), 2999970, c)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(m)
	fmt.Println(string(b))

	if m.Credits != -95175 {
		t.FailNow()
	}
}
