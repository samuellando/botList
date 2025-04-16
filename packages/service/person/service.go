package person

import (
    "crypto/ed25519"
	"crypto/rand"
    "encoding/base64"
	"encoding/json"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/action"
	"fedilist/packages/model/list"
	"fedilist/packages/model/person"
	"fedilist/packages/util"
	"fmt"
    "io"
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

func Create(store PersonStore, q chan []byte) PersonService {
	return PersonService{
		store:        store,
		messageQueue: q,
	}
}

func (ls PersonService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	endpoint := req.PathValue("endpoint")
    if endpoint == "" && req.Method == "POST"  {
        bodyBytes, err := io.ReadAll(req.Body)
        if err != nil {
            panic(err)
        }
        defer req.Body.Close()
        data := make(map[string]any)
        json.Unmarshal(bodyBytes, &data)
        p, privateKey, err := ls.Create(func(pv *person.PersonValues) {
            pv.Name = data["name"].(string)
        })
        out, err := json.MarshalIndent(map[string]any{
            "id": p.Id(),
            "privateKey": privateKey,
        }, "", "    ")
        if err != nil {
            panic(err)
        }
        w.Write(out)
        return
    }

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

func (s PersonService) Create(fs ...func(*person.PersonValues)) (person.Person, string, error) {
    publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
    if err != nil {
        panic(err)
    }
    fs = append(fs, func(pv *person.PersonValues) {
        pv.Key = base64.StdEncoding.EncodeToString(publicKey)
    })
    p, err := s.store.Insert(person.CreatePerson(fs...))
    if err != nil {
        return p, "", err
    }
    return p, base64.StdEncoding.EncodeToString(privateKey.Seed()), nil
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
