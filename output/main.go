package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Points map[string]string `json:"points"`
}

func readConfig() config {
	c := config{}
	b, err := os.ReadFile(`config.json`)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}
	return c
}

func main() {
	_ = readConfig()
}
