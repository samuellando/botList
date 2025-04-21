package main

import (
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

	ls := listService.Create(listStore.CreateStore(serverUrl()+"/list/"), messages)
	ps := personService.Create(personStore.CreateStore(serverUrl()+"/user/"), messages, ls)
	rs := runnerService.Create(serverUrl()+"/runner", messages)

	http.Handle("/list/{id}/{endpoint...}", ls)

	http.Handle("/user", ps)

	http.Handle("/user/{id}/{endpoint...}", ps)

	http.Handle("/runner/{endpoint...}", rs)

	http.ListenAndServe(":"+args[0], nil)
}
