package profile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGetCredit(t *testing.T) {
	m, err := GetCredit(80321, &http.Client{})
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(m)
	fmt.Println(string(b))
}
