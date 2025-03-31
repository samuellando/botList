package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

func serverUrl() string {
	args := os.Args[1:]
	return fmt.Sprintf("http://localhost:%s", args[0])
}

func getId(typ, id string) string {
	return fmt.Sprintf("%s/%s/%s", serverUrl(), typ, id)
}

func getReqObjectReference(req *http.Request) (string, error) {
	baseURL, err := url.Parse(serverUrl())
	if err != nil {
		return "", err
	}
	listReference, err := url.Parse(path.Dir(req.URL.Path))
	if err != nil {
		return "", err
	}
	objectId := baseURL.ResolveReference(listReference).String()
	return objectId, nil
}

func getReqActivity[T any](req *http.Request) (Activity[T], error) {
	bodyBytes, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return Activity[T]{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	var activity Activity[T]
	err = decoder.Decode(&activity)
	if err != nil {
		return Activity[T]{}, err
	}
	return activity, nil
}

func main() {
	args := os.Args[1:]

    p := CreatePerson("sam", "")
    p.Lists.AddHook(Hook{Type: "onAdd", Integration: "Add", Runner: serverUrl()+"/runner"})
    p.Lists.AddHook(Hook{Type: "onAdd", Integration: "Add", Runner: serverUrl()+"/runner"})
    p.Lists.AddHook(Hook{Type: "cron",  Integration: "Tag", Runner: serverUrl()+"/runner", Cron: "0,30 * * * *"})

	messages := make(chan string)
	go ProcessMessages(messages)

	http.HandleFunc("/user/{username}/outbox", func(w http.ResponseWriter, req *http.Request) {
        userId, err := getReqObjectReference(req)
        if err != nil {
            panic(err)
        }
        activity, err := getReqActivity[List](req)
		if userId != activity.Actor {
			http.Error(w, "Can't post to another user's outbox", 400)
			return
		}
		switch req.Method {
		case "POST":
            jsonB, err := json.Marshal(activity)
            if err != nil {
                panic(err)
            }
			messages <- string(jsonB)
			w.WriteHeader(202)
			return
		default:
			http.Error(w, "Only POST is supported to outbox", 400)
			return
		}
	})

	http.HandleFunc("/user/{username}/inbox", func(w http.ResponseWriter, req *http.Request) {
        userId, err := getReqObjectReference(req)
        if err != nil {
            panic(err)
        }
		p := GetPersonById(userId)
		switch req.Method {
		case "GET":
			en := json.NewEncoder(w)
			err := en.Encode(p.Inbox)
			if err != nil {
				panic(err)
			}
			return
		case "POST":
            activity, err := getReqActivity[List](req)
            if err != nil {
                panic(err)
            }
			if activity.Result.Type != "" {
				p.AddToInbox(activity)
				w.WriteHeader(202)
			}
			return
		default:
			http.Error(w, "Only POST is supported to outbox", 400)
			return
		}
	})

	http.HandleFunc("/list/{id}/inbox", func(w http.ResponseWriter, req *http.Request) {
		listId, err := getReqObjectReference(req)
		if err != nil {
			panic(err)
		}
		activity, err := getReqActivity[List](req)
		if err != nil {
			panic(err)
		}
        fmt.Println("OK", activity, listId)
		if listId != activity.Target {
			http.Error(w, "Activity target does not match.", 400)
			return
		}
		l := GetListById(activity.Target)
		switch req.Method {
		case "POST":
			go func() {
				result := l.ProcessActivity(activity)
				responseBuilder := new(strings.Builder)
				en := json.NewEncoder(responseBuilder)
				err := en.Encode(result)
				if err != nil {
					panic(err)
				}
				messages <- responseBuilder.String()
			}()
			w.WriteHeader(202)
			return
		default:
			http.Error(w, "Only POST is supported to outbox", 400)
			return
		}
	})

	http.HandleFunc("/runner/inbox", func(w http.ResponseWriter, req *http.Request) {
		activity, err := getReqActivity[Activity[List]](req)
		if err != nil {
			panic(err)
		}
		r := CreateRunner(serverUrl()+"/runner")
		switch req.Method {
		case "POST":
			go func() {
				r.ProcessActivity(activity)
				// respB, err := json.Marshal(result)
				// if err != nil {
				// 	panic(err)
				// }
				// messages <- string(respB)
			}()
			w.WriteHeader(202)
			return
		default:
			http.Error(w, "Only POST is supported to outbox", 400)
			return
		}
	})

	http.HandleFunc("/list/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			panic(err)
		}
		PrintList(w, GetListById(fmt.Sprintf("%s/list/%d", serverUrl(), id)))
	})

	http.ListenAndServe(":"+args[0], nil)
}

func PrintList(w io.Writer, l List) {
	fmt.Fprintln(w, "----------------")
	fmt.Fprintln(w, l.Name)
	fmt.Fprintln(w, "----------------")
	for _, e := range l.Items {
		fmt.Fprintln(w, "\t-", e.Name)
	}
}
