package main

import (
	"fedilist/packages/parser/action"
	"fedilist/packages/parser/jsonld"
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/person"
	"fedilist/packages/parser/result"
    "strconv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func serverUrl() string {
	args := os.Args[1:]
	return fmt.Sprintf("http://localhost:%s", args[0])
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

func getBodyJsonld(req *http.Request) (map[string]any, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	data, err := jsonld.Expand(bodyBytes)
	return data, err
}

func main() {
	args := os.Args[1:]

	p := person.CreatePerson(serverUrl(), "sam", "")
	l := list.CreateList(serverUrl(), "Sam's list", "")
	p.AddList(l)

	messages := make(chan []byte, 100)
	go ProcessMessages(messages)

	http.HandleFunc("/user/{username}/outbox", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Only POST is supported to outbox", 400)
			return
		}
		userId, err := getReqObjectReference(req)
		if err != nil {
			panic(err)
		}
		json, err := getBodyJsonld(req)
		if err != nil {
			panic(err)
		}
		anyAct, err := action.Parse(json)
		if err != nil {
			panic(err)
		}
		if userId != anyAct.Agent().Id {
			http.Error(w, "Can't post to another user's outbox", 400)
			return
		}
		switch act := anyAct.(type) {
		case action.Create:
		case action.Action:
			jsonB, err := act.Serialize()
			if err != nil {
				panic(err)
			}
			messages <- jsonB
			w.WriteHeader(202)
			return
		default:
			http.Error(w, "Can only post actions", 400)
			return
		}
	})

	http.HandleFunc("/list/{id}/inbox", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Only POST is supported to outbox", 400)
			return
		}
		listId, err := getReqObjectReference(req)
		if err != nil {
			panic(err)
		}
		json, err := getBodyJsonld(req)
		if err != nil {
			panic(err)
		}
		anyAct, err := action.Parse(json)
		if err != nil {
			panic(err)
		}
		if listId != *anyAct.TargetId() {
			http.Error(w, "Can't post to lists user's inbox", 400)
			return
		}
        switch act := anyAct.(type) {
		case action.Append:
            l := list.GetListById(listId)
            if act.Object().Name == nil {
                panic("NO NAME")
            }
            e := list.CreateList(serverUrl(), *act.Object().Name, "")
            l.Append(e)
            act.AddResult(result.Create("201", "Appended"))
			jsonB, err := act.Serialize()
			if err != nil {
				panic(err)
			}
			messages <- jsonB
            fmt.Println("OK")
			w.WriteHeader(202)
			return
		default:
			http.Error(w, "Unsupported action", 400)
			return
		}
	})

	http.HandleFunc("/list/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			panic(err)
		}
		PrintList(w, list.GetListById(fmt.Sprintf("%s/list/%d", serverUrl(), id)))
	})

	http.HandleFunc("/user/{username}/inbox", func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("Got response")
        w.WriteHeader(202)
	})

	http.ListenAndServe(":"+args[0], nil)

	//
	// http.HandleFunc("/runner/inbox", func(w http.ResponseWriter, req *http.Request) {
	// 	activity, err := getReqActivity[Activity[List]](req)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	r := CreateRunner(serverUrl() + "/runner")
	// 	switch req.Method {
	// 	case "POST":
	// 		go func() {
	// 			r.ProcessActivity(activity)
	// 			// respB, err := json.Marshal(result)
	// 			// if err != nil {
	// 			// 	panic(err)
	// 			// }
	// 			// messages <- string(respB)
	// 		}()
	// 		w.WriteHeader(202)
	// 		return
	// 	default:
	// 		http.Error(w, "Only POST is supported to outbox", 400)
	// 		return
	// 	}
	// })
	//

}

func PrintList(w io.Writer, l list.ItemList) {
	fmt.Fprintln(w, "----------------")
	fmt.Fprintln(w, *l.Name)
	fmt.Fprintln(w, "----------------")
	for _, e := range l.ItemListElement {
		fmt.Fprintln(w, "\t-", *e.Name)
	}
}
