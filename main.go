package main

import (
	"fedilist/packages/parser/action"
	"fedilist/packages/parser/jsonld"
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/person"
	"fedilist/packages/parser/result"
	listService "fedilist/packages/service/list"
	listStore "fedilist/packages/store/list"
	"fedilist/packages/util"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
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

	p := person.CreatePerson(serverUrl(), "sam", "")
	l, _ := ls.Create(func(ilv *list.ItemListValues) {
		name := "Sam's list"
		ilv.Name = &name
		r, err := list.CreateRunner(func(rv *list.RunnerValues) {
			rv.Id = serverUrl() + "/runner"
			rv.Name = "Default Runner"
			rv.Inbox = serverUrl() + "/runner/inbox"
		})
		if err != nil {
			panic(err)
		}
		h, err := list.CreateActionHook(func(ahv *list.ActionHookValues) {
			ahv.Runner = r
			ahv.RunnerAction = "CopyTo"
			ahv.RunnerActionConfig = serverUrl() + "/list/1"
			ahv.OnActionType = []string{"AppendAction"}
		})
		ch, err := list.CreateCronHook(func(ahv *list.CronHookValues) {
			ahv.Runner = r
			ahv.RunnerAction = "Print"
			ahv.RunnerActionConfig = "Hello World"
			ahv.CronTab = "0,30 * * * *"
		})
		ilv.Hooks = []list.Hook{h, ch}
	})
	ls.Create(func(ilv *list.ItemListValues) {
		name := "Target list"
		ilv.Name = &name
	})
	p.AddList(l)

	http.Handle("/list/{id}/{endpoint...}", ls)

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

	http.HandleFunc("/runner/inbox", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Got Request!!")
		json, err := util.GetBodyJsonld(req)
		if err != nil {
			panic(err)
		}
		act, err := action.Parse(json)
		if err != nil {
			panic(err)
		}
		switch ea := act.(type) {
		case action.Execute:
			switch ea.RunnerAction() {
			case "Print":
				fmt.Println("RUNNER:", ea.RunnerActionConfig())
				w.WriteHeader(202)
			case "CopyTo":
				ca := action.CreateAppend(func(av *action.AppendValues) {
					oa := ea.Object().(action.Append)
					tl, err := ls.GetById(ea.RunnerActionConfig())
					if err != nil {
						panic(err)
					}
					av.Agent = ea.Agent()
					av.Object = oa.Object()
					av.StartTime = time.Now()
					av.TargetCollection = tl
				})
				messages <- jsonld.Marshal(ca)

				act = act.WithResult(result.Create("201", "RAN HOOK"))
				messages <- jsonld.Marshal(act)
				w.WriteHeader(202)
			default:
				panic("Not a runner action")
			}
		default:
			panic("ONLY EXECUTE")
		}
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
