package list

import (
	"fedilist/packages/model/action"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/list"
	"fedilist/packages/model/hook"
	"fedilist/packages/model/result"
	"fedilist/packages/service/cron"
	"fedilist/packages/util"
	"net/http"
	"time"
)

type ListStore interface {
	GetById(string) (list.ItemList, error)
	GetByPartialId(string) (list.ItemList, error)
	Insert(list.ItemList) (list.ItemList, error)
	Append(list.ItemList, list.ItemList) (list.ItemList, error)
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
        cronService: cron.Create(q),
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

func (s ListService) Create(fs ...func(*list.ItemListValues)) (list.ItemList, error) {
    l, err := s.store.Insert(list.Create(fs...))
    if err != nil {
        return l, err
    }
    for _, h := range l.Hooks() {
        switch ch := h.(type) {
        case hook.CronHook:
            ea := action.CreateExecute(func(ev *action.ExecuteValues) {
                ev.StartTime = time.Now()
                ev.TargetRunner = ch.Runner()
                ev.RunnerAction = ch.RunnerAction()
                ev.RunnerActionConfig = ch.RunnerActionConfig()
            })
            b := jsonld.MarshalIndent(ea)
            s.cronService.AddJob(cron.CronJob{
                Crontab: ch.CronTab(),
                Message: b,
            })
        }
    }
    return l, nil
}

func (ls ListService) Append(w http.ResponseWriter, act action.Append) {
	if act.Object().Name() == nil {
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

	_, err = ls.store.Append(target, e)
	if err != nil {
		panic(err)
	}
	jsonB := jsonld.Marshal(act.WithResult(result.Create("201", "Appended")))

	ls.messageQueue <- jsonB
	w.WriteHeader(202)
    l, err := ls.GetById(*target.Id())
    if err != nil {
        panic(err)
    }
    for _, h := range l.Hooks() {
        switch h.(type) {
        case hook.ActionHook:
            ea := action.CreateExecute(func(ev *action.ExecuteValues) {
                ev.Agent= act.Agent()
                ev.Object = act
                ev.StartTime = time.Now()
                ev.TargetRunner = h.Runner()
                ev.RunnerAction = h.RunnerAction()
                ev.RunnerActionConfig = h.RunnerActionConfig()
            })
            b := jsonld.MarshalIndent(ea)
            ls.messageQueue <- b
        }
    }
}
