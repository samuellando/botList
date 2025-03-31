package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Runner struct {
    Id string
}

func CreateRunner(id string) Runner {
    return Runner{Id: id}
}

func (r Runner) ProcessActivity(activity Activity[Activity[List]]) Activity[Activity[List]] {
    var req Activity[List]
    switch activity.Integration {
    case "Add":
        req = Activity[List]{
            Context: "https://www.w3.org/ns/activitystreams",
            Type: "Add",
            Actor: r.Id,
            Object: List{
                Name: "Runner",
            },
            Target: activity.Actor,
        }
        reqB, err := json.Marshal(req)
        if err != nil {
            panic(err)
        }
        resp, err := http.Post(activity.Actor+"/inbox", "text/json", bytes.NewReader(reqB))
        if err != nil {
            panic(err)
        }
        if resp.StatusCode != 202 {
            panic("Server responded with a bad error code")
        }
        activity.Result = CreateResult(202, "Submitted")
        return activity
    case "Tag":
        req = Activity[List]{
            Context: "https://www.w3.org/ns/activitystreams",
            Type: "Update",
            Actor: r.Id,
            Object: List{
                Name: "Title Runner",
            },
            Target: activity.Actor,
        }
        reqB, err := json.Marshal(req)
        if err != nil {
            panic(err)
        }
        resp, err := http.Post(activity.Actor+"/inbox", "text/json", bytes.NewReader(reqB))
        if err != nil {
            panic(err)
        }
        if resp.StatusCode != 202 {
            panic("Server responded with a bad error code")
        }
        activity.Result = CreateResult(202, "Submitted")
        return activity
    }
    activity.Result = CreateResult(404, "Integration not found")
    return activity
}
