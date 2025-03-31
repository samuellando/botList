package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/robfig/cron"
)

type List struct {
	Id           string
	Name         string
	Summary      *string
	AttributedTo []string
	Editors      []string
	Viewers      []string
	Hooks        []Hook
	Tags         []string
	TotalItems   int
	// Can either have a reference to itself, or an actual set of sublists
	// Allowing for sublists to not load all sub elements
	First string
	Items []List
}

type Hook struct {
	Type        string
	Runner      string
	Integration string
	Cron        string
}

var LISTS = make(map[string]*List)

var ID = 0

func GetListById(id string) List {
	return *LISTS[id]
}

func CreateList(name, attributed string) List {
	l := List{
		Id:           getId("list", strconv.Itoa(ID)),
		Name:         name,
		AttributedTo: []string{attributed},
		Items:        make([]List, 0),
	}
	LISTS[l.Id] = &l
	ID++
	return l
}

func (l *List) AddHook(hook Hook) {
	l.Hooks = append(l.Hooks, hook)
	LISTS[l.Id].Hooks = append(LISTS[l.Id].Hooks, hook)
	if hook.Type == "cron" {
		go func() {
			specParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional)
			sched, err := specParser.Parse(hook.Cron)
			if err != nil {
				panic(err)
			}
			next := sched.Next(time.Now())
			for {
				if next.Sub(time.Now()) <= 0 {
					fmt.Println("Running")
					req := Activity[Activity[List]]{
						Context: "https://www.w3.org/ns/activitystreams",
						Type:    "RunHook",
                        Integration: hook.Integration,
                        Target: hook.Runner,
						Actor:   l.Id,
					}
					reqB, err := json.Marshal(req)
					if err != nil {
						panic(err)
					}
					resp, err := http.Post(hook.Runner+"/inbox", "text/json", bytes.NewReader(reqB))
					if err != nil {
						panic(err)
					}
					if resp.StatusCode != 202 {
						panic("Server responded with a bad error code")
					}
					next = sched.Next(time.Now())
				} else {
					fmt.Println("Sleeping")
					time.Sleep(next.Sub(time.Now()))
				}
			}
		}()
	}
}

func (l *List) ProcessActivity(activity Activity[List]) Activity[List] {
	switch activity.Type {
	case "Add":
        // Check is editor
		e := CreateList(activity.Object.Name, activity.Actor)
		l.Add(e)
		activity.Object.Id = e.Id
		l.runHooks("onAdd", activity)
	case "Update":
        // Check is editor
		l.Name = activity.Object.Name
		l.Summary = activity.Object.Summary
		l.Tags = activity.Object.Tags
        LISTS[l.Id] = l
		l.runHooks("onUpdate", activity)
	case "Clear":
        // Check is editor
        l.Items = make([]List, 0)
        l.TotalItems = 0
        LISTS[l.Id] = l
		l.runHooks("onClear", activity)
	}
	// Return the result.
	activity.Result = CreateResult(201, "Added")
	return activity
}

func (l *List) Add(e List) List {
	LISTS[l.Id].Items = append(LISTS[l.Id].Items, e)
	l.Items = append(l.Items, e)
	return *l
}

func (l *List) runHooks(typ string, activity Activity[List]) {
	req := Activity[Activity[List]]{
		Context: "https://www.w3.org/ns/activitystreams",
		Type:    "RunHook",
		Actor:   l.Id,
		Object:  activity,
	}
	for _, hook := range l.Hooks {
		if hook.Type != typ {
			continue
		}
		if hook.Runner == activity.Actor {
			// Don't hook on self
			continue
		}
		req.Target = hook.Runner
		req.Integration = hook.Integration
		reqB, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}
		resp, err := http.Post(hook.Runner+"/inbox", "text/json", bytes.NewReader(reqB))
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 202 {
			panic("Server responded with a bad error code")
		}
	}
}
