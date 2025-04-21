package list

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/action"
	"fedilist/packages/model/hook"
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"fedilist/packages/service/cron"
	"fedilist/packages/util"
	"fmt"
	"net/http"
	"time"
)

type ListStore interface {
	GetById(string) (list.ItemList, error)
	GetByPartialId(string) (list.ItemList, error)
	Insert(list.ItemList) (list.ItemList, error)
	Append(list.ItemList, list.ItemList) (list.ItemList, error)
	GetKey(list.ItemList) ([]byte, error)
	StoreKey(list.ItemList, []byte) error
}

type ListService struct {
	store        ListStore
	messageQueue chan []byte
	cronService  cron.CronService
}

func Create(store ListStore, q chan []byte) ListService {
	return ListService{
		store:        store,
		messageQueue: q,
		cronService:  cron.Create(q),
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
		s := jsonld.MarshalIndent(list)
		w.Write(s)
	case "inbox":
		if req.Method != "POST" {
			http.Error(w, "Only POST is supported to inboxes", 400)
			return
		}
		json, err := util.GetBodyJsonld(req.Body)
		if err != nil {
			panic(err)
		}
		anyAct, err := action.Parse(json)
		if list.Id() != *anyAct.TargetId() && (anyAct.Result() != nil && anyAct.Agent().Id() != list.Id()) {
			http.Error(w, "Action target does not match request URL", 400)
			return
		}
		var resp *http.Response
		if anyAct.Result() == nil {
			resp, err = http.Get(anyAct.Agent().Id())
		} else {
			resp, err = http.Get(*anyAct.TargetId())
		}
		personJson, err := util.GetBodyJsonld(resp.Body)
		if err != nil {
			panic(err)
		}
		agent, err := action.ParseAgent(personJson)
		if err != nil {
			panic(err)
		}
		valid, err := util.VerifySignature(anyAct, agent.Key())
		if err != nil {
			panic(err)
		}
		if !valid {
			http.Error(w, "Could not verify message signature", 403)
			return
		}
		if anyAct.Result() != nil {
			fmt.Println("LIST GOT RESULT")
			w.WriteHeader(202)
			return
		}
		switch act := anyAct.(type) {
		case action.Append:
			ls.Append(w, act)
		default:
			http.Error(w, "Unsupported action", 400)
			return
		}
	}
}

func (s ListService) GetById(id string) (list.ItemList, error) {
	return s.store.GetById(id)
}

func (s ListService) Create(fs ...func(*list.ItemListValues)) (list.ItemList, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	fs = append(fs, func(ilv *list.ItemListValues) {
		ilv.Key = base64.StdEncoding.EncodeToString(publicKey)
	})
	l, err := s.store.Insert(list.Create(fs...))
	if err != nil {
		return l, err
	}
	err = s.store.StoreKey(l, privateKey.Seed())
	if err != nil {
		return l, err
	}
	for _, h := range l.Hooks() {
		switch ch := h.(type) {
		case hook.CronHook:
			ea := action.CreateExecute(func(ev *action.ExecuteValues) {
				ev.Agent = l
				ev.StartTime = time.Now()
				ev.TargetRunner = ch.Runner()
				ev.RunnerAction = ch.RunnerAction()
				ev.RunnerActionConfig = ch.RunnerActionConfig()
			})
			pk, err := s.store.GetKey(l)
			if err != nil {
				panic(err)
			}
			b := jsonld.MarshalIndent(util.Sign(ea, pk))
			s.cronService.AddJob(cron.CronJob{
				Crontab: ch.CronTab(),
				Message: b,
			})
		}
	}
	return l, nil
}

func (ls ListService) Append(w http.ResponseWriter, act action.Append) {
	if act.Object().Name() == "" {
		http.Error(w, "Action Object requires a name", 400)
		return
	}
	obj := act.Object()
	e, err := ls.Create(func(ilv *list.ItemListValues) {
		ilv.Name = obj.Name()
		ilv.Description = obj.Description()
	})
	if err != nil {
		panic(err)
	}
	target := act.TargetCollection()
	target.Append(e)

	l, err := ls.GetById(target.Id())
	if err != nil {
		panic(err)
	}
	pk, err := ls.store.GetKey(l)
	if err != nil {
		panic(err)
	}

	_, err = ls.store.Append(target, e)
	if err != nil {
		panic(err)
	}
	res := act.WithResult(result.Create("201", "Appended"))
	jsonB := jsonld.Marshal(util.Sign(res, pk))
	ls.messageQueue <- jsonB
	w.WriteHeader(202)

	for _, h := range l.Hooks() {
		switch h.(type) {
		case hook.ActionHook:
			ea := action.CreateExecute(func(ev *action.ExecuteValues) {
				ev.Agent = l
				ev.Object = act
				ev.StartTime = time.Now()
				ev.TargetRunner = h.Runner()
				ev.RunnerAction = h.RunnerAction()
				ev.RunnerActionConfig = h.RunnerActionConfig()
			})
			b := jsonld.MarshalIndent(util.Sign(ea, pk))
			ls.messageQueue <- b
		}
	}
}
