package customhttp

import (
	"encoding/json"
	"log"
)

type resp struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func ApiResp(value interface{}) ([]byte, error) {
	r := resp{Data: value, Msg: "OK"}
	bytes, err := json.Marshal(r)
	if err != nil {
		r.Msg = "error while marshalling JSON"
		log.Printf("error while marshalling JSON: %s", err)
	}
	return bytes, err
}
