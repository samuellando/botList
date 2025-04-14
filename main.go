package main

import (
	"fedilist/packages/model/hook"
	"fedilist/packages/model/list"
	"fedilist/packages/model/person"
	listService "fedilist/packages/service/list"
	personService "fedilist/packages/service/person"
	"fedilist/packages/service/router"
	runnerService "fedilist/packages/service/runner"
	listStore "fedilist/packages/store/list"
	personStore "fedilist/packages/store/person"
	"fmt"
	"net/http"
	"os"
)

func serverUrl() string {
	args := os.Args[1:]
	return fmt.Sprintf("http://localhost:%s", args[0])
}

func main() {
	args := os.Args[1:]

	messages := make(chan []byte, 100)
	go router.ProcessMessages(messages)

	ls := listService.CreateService(listStore.CreateStore(serverUrl()+"/list/"), messages)
	ps := personService.CreateService(personStore.CreateStore(serverUrl()+"/user/"), messages)
	rs := runnerService.Create(serverUrl()+"/runner", messages)

	p, _ := ps.Create(func(pv *person.PersonValues) {
		pv.Name = "Sam"
	})

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

	ps.AddList(p, l)

	http.Handle("/list/{id}/{endpoint...}", ls)

	http.Handle("/user/{id}/{endpoint...}", ps)

	http.Handle("/runner/{endpoint...}", rs)

	http.ListenAndServe(":"+args[0], nil)
}
