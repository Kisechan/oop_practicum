package service

import "encoding/json"

var ReqHandler = map[string]func(Request) (json.RawMessage, error){}
