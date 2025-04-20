package router

import (
	"bytes"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/action"
	"fmt"
	"io"
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
            to = act.Agent().Id()
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
			txt, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			panic("Server responded with a bad error code: "+string(txt))
        }
        continue
	}
}
