package list

import (
	"fedilist/packages/jsonld"
	"fedilist/packages/model/action"
	"fedilist/packages/model/hook"
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"fedilist/packages/service/cron"
	"fedilist/packages/util"
	"net/http"
	"slices"
	"time"
)

type ListStore interface {
	GetById(string) (list.ItemList, error)
	GetByPartialId(string) (list.ItemList, error)
	Insert(list.ItemList) (list.ItemList, error)
	Append(list.ItemList, list.ItemList) (list.ItemList, error)
	Remove(list.ItemList, list.ItemList) (list.ItemList, error)
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
		ls.handleInbox(list, w, req)
	}
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

	l, err := ls.store.GetById(target.Id())
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

func (ls ListService) Remove(w http.ResponseWriter, act action.Remove) {
	obj := act.Object()
	target := act.TargetCollection()

	l, err := ls.store.GetById(target.Id())
	if err != nil {
		panic(err)
	}
	pk, err := ls.store.GetKey(l)
	if err != nil {
		panic(err)
	}

	_, err = ls.store.Remove(target, obj)
	if err != nil {
		panic(err)
	}
	res := act.WithResult(result.Create("201", "Appended"))
	jsonB := jsonld.Marshal(util.Sign(res, pk))
	ls.messageQueue <- jsonB
	w.WriteHeader(202)

	for _, h := range l.Hooks() {
		switch ah := h.(type) {
		case hook.ActionHook:
			if !slices.Contains(ah.OnActionType(), "Remove") {
				continue
			}
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
