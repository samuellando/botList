package main

import (
	"encoding/json"
	"net/http"
	"strings"
    "fmt"
)

func ProcessMessages(q chan string) {
	for {
		msg := <-q
        decoder := json.NewDecoder(strings.NewReader(msg))
        var activity Activity[List]
        err := decoder.Decode(&activity)
        if err != nil {
            panic(err)
        }
        var to string
        if activity.Result.Type == "" {
            to = activity.Target
        } else {
            to = activity.Actor
        }
        fmt.Println(">", to+"/inbox")
        resp, err := http.Post(
            to+"/inbox", "text/json",
            strings.NewReader(msg),
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
