package main

import (
	"fedilist/packages/model/action"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/list"
	"fedilist/packages/model/hook"
	"fedilist/packages/model/person"
	listService "fedilist/packages/service/list"
	runnerService "fedilist/packages/service/runner"
	listStore "fedilist/packages/store/list"
	"fedilist/packages/util"
	"fmt"
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

func main() {
	args := os.Args[1:]

	messages := make(chan []byte, 100)
	go ProcessMessages(messages)


	ls := listService.CreateService(listStore.CreateStore(serverUrl()+"/list/"), messages)
    rs := runnerService.Create(serverUrl()+"/runner", messages)

	p := person.CreatePerson(serverUrl(), "sam", "")
	l, _ := ls.Create(func(ilv *list.ItemListValues) {
		name := "Sam's list"
		ilv.Name = &name
		h, err := hook.CreateActionHook(func(ahv *hook.ActionHookValues) {
			ahv.Runner = rs.Runner()
			ahv.RunnerAction = "CopyTo"
			ahv.RunnerActionConfig = serverUrl() + "/list/1"
			ahv.OnActionType = []string{"AppendAction"}
		})
        if err != nil {
            panic(err)
        }
		ch, err := hook.CreateCronHook(func(ahv *hook.CronHookValues) {
			ahv.Runner = rs.Runner()
			ahv.RunnerAction = "Print"
			ahv.RunnerActionConfig = "Hello World"
			ahv.CronTab = "0,30 * * * *"
		})
        if err != nil {
            panic(err)
        }
		ilv.Hooks = []hook.Hook{h, ch}
	})
	ls.Create(func(ilv *list.ItemListValues) {
		name := "Target list"
		ilv.Name = &name
	})
	p.AddList(l)

	http.Handle("/list/{id}/{endpoint...}", ls)

	http.Handle("/runner/{endpoint...}", rs)

	http.HandleFunc("/user/{username}/outbox", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Only POST is supported to outbox", 400)
			return
		}
		userId, err := getReqObjectReference(req)
		if err != nil {
			panic(err)
		}
		json, err := util.GetBodyJsonld(req)
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
			jsonB := jsonld.Marshal(act)
			messages <- jsonB
			w.WriteHeader(202)
			return
		default:
			http.Error(w, "Can only post actions", 400)
			return
		}
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
