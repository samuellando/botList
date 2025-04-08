package list

import (
	"fedilist/packages/parser/action"
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/result"
	"fedilist/packages/util"
	"fmt"
	"net/http"
)

type ListStore interface {
	GetById(string) (list.ItemList, error)
	GetByPartialId(string) (list.ItemList, error)
	Insert(list.ItemList) (list.ItemList, error)
	Append(list.ItemList, list.ItemList) (list.ItemList, error)
}

type ListService struct {
	store ListStore
    messageQueue chan []byte
}

func CreateService(store ListStore, q chan []byte) ListService {
	return ListService{
		store: store,
        messageQueue: q,
	}
}

func (ls ListService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	endpoint := req.PathValue("endpoint")
	list, err := ls.store.GetByPartialId(id)
	if err != nil {
		panic(err)
	}
	switch endpoint {
	case "":
        fmt.Println(list)
		w.Write(list.Serialize())
	case "inbox":
		if req.Method != "POST" {
			http.Error(w, "Only POST is supported to inboxes", 400)
			return
		}
		json, err := util.GetBodyJsonld(req)
		if err != nil {
			panic(err)
		}
		anyAct, err := action.Parse(json)
		if *list.Id() != *anyAct.TargetId() {
			http.Error(w, "Action target does not match request URL", 400)
			return
		}
		switch act := anyAct.(type) {
		case action.Append:
            ls.Append(w, act)
            if err != nil {
                panic(err)
            }

		default:
			http.Error(w, "Unsupported action", 400)
			return
		}
	}
}

func (s ListService) GetById(id string) (list.ItemList, error) {
	return s.store.GetById(id)
}

func (s ListService) Create(l list.ItemList) (list.ItemList, error) {
	return s.store.Insert(list.Create(func(ilv *list.ItemListValues) {
		ilv.Name = l.Name()
		ilv.Description = l.Description()
	}))
}


func (ls ListService) Append(w http.ResponseWriter, act action.Append){
	if act.Object().Name() == nil {
		http.Error(w, "Action Object requires a name", 400)
        return
	}
    obj := act.Object()
	e, err  := ls.Create(obj)
    if err != nil {
        panic(err)
    }
    target := act.TargetCollection()
	target.Append(e)

    _, err = ls.store.Append(target, e)
    if err != nil {
        panic(err)
    }

	act.AddResult(result.Create("201", "Appended"))
	jsonB := act.Serialize()
	ls.messageQueue <- jsonB
	w.WriteHeader(202)
}
