package main

import (
	"fedilist/packages/parser/action"
	"fedilist/packages/parser/jsonld"
	"fmt"
	"bytes"
    "net/http"
)

func ProcessMessages(q chan []byte) {
	for {
		msg := <-q
        fmt.Println("GOT MESSAGE")
        fmt.Println(string(msg))
        data, err := jsonld.Expand([]byte(msg))
        if err != nil {
            panic(err)
        }
        act, err := action.Parse(data)
        if err != nil {
            panic(err)
        }
        var to string
        if act.Result() == nil {
            to = *act.TargetId()
        } else {
            to = act.Agent().Id
        }
        fmt.Println(">", to+"/inbox")
        resp, err := http.Post(
            to+"/inbox", "text/json",
            bytes.NewReader(msg),
        )
        if err != nil {
            panic(err)
        }
        if resp.StatusCode != 202 {
            panic("Server responded with a bad error code")
        }
        continue
	}
}
