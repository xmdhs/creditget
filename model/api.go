package model

type ApiRep[V any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data V      `json:"data"`
}

type ApiErr int

const (
	ApiOk ApiErr = iota
	ApiDateBaseFail
	ApiErrInput
)
