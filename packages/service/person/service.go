package person

import (
	"fedilist/packages/jsonld"
	"fedilist/packages/model/action"
	"fedilist/packages/model/list"
	"fedilist/packages/model/person"
	"fedilist/packages/util"
	"fmt"
	"net/http"
)

type PersonStore interface {
	GetById(id string) (person.Person, error)
	GetByPartialId(pid string) (person.Person, error)
	Insert(person.Person) (person.Person, error)
	AddList(person.Person, list.ItemList) (person.Person, error)
}

type PersonService struct {
	store        PersonStore
	messageQueue chan []byte
}

func CreateService(store PersonStore, q chan []byte) PersonService {
	return PersonService{
		store:        store,
		messageQueue: q,
	}
}

func (ls PersonService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	endpoint := req.PathValue("endpoint")
	person, err := ls.store.GetByPartialId(id)
	if err != nil {
		panic(err)
	}
	switch endpoint {
	case "":
		s := jsonld.MarshalIndent(person)
		w.Write(s)
	case "outbox":
		ls.handleOutbox(person, w, req)
	case "inbox":
		ls.handleInbox(person, w, req)
	}
}

func (s PersonService) Create(fs ...func(*person.PersonValues)) (person.Person, error) {
    return s.store.Insert(person.CreatePerson(fs...))
}

func (s PersonService) AddList(p person.Person, l list.ItemList) (person.Person, error) {
    return s.store.AddList(p, l)
}

func (ls PersonService) handleOutbox(p person.Person, w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Only POST is supported to outbox", 400)
		return
	}
	json, err := util.GetBodyJsonld(req)
	if err != nil {
		panic(err)
	}
	anyAct, err := action.Parse(json)
	if err != nil {
		panic(err)
	}
	if p.Id() != anyAct.Agent().Id() {
		http.Error(w, "Can't post to another user's outbox", 400)
		return
	}
	switch act := anyAct.(type) {
	case action.Create:
	case action.Action:
		jsonB := jsonld.Marshal(act)
		ls.messageQueue <- jsonB
		w.WriteHeader(202)
		return
	default:
		http.Error(w, "Can only post actions", 400)
		return
	}
}

func (ls PersonService) handleInbox(p person.Person, w http.ResponseWriter, req *http.Request) {
	fmt.Println("Got response")
	w.WriteHeader(202)
}
