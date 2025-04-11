package runner

import (
	"fedilist/packages/model/action"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/list"
	"fedilist/packages/model/runner"
	"fedilist/packages/model/result"
	"fedilist/packages/util"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	messageQueue chan []byte
	runner       runner.Runner
}

func Create(id string, q chan []byte) Service {
	r, err := runner.Create(func(rv *runner.RunnerValues) {
		rv.Id = id
		rv.Name = "Default Runner"
		rv.Inbox = id + "/inbox"
		rv.Service = []runner.Service{
			{
				Name:   "CopyTo",
				Schema: `{"example": "ok"}`,
			},
			{
				Name:   "CloneList",
				Schema: `{"example": "ok"}`,
			},
		}
	})
	if err != nil {
		panic(err)
	}
	return Service{
		messageQueue: q,
		runner:       r,
	}
}

func (s Service) Runner() runner.Runner {
	return s.runner
}

func (rs Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	endpoint := req.PathValue("endpoint")
	switch endpoint {
	case "":
		w.Write(jsonld.MarshalIndent(rs.runner))
	case "inbox":
		if req.Method != "POST" {
			http.Error(w, "Only POST is supported to inboxes", 400)
			return
		}
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
					tl := list.Create(func(ilv *list.ItemListValues) {
						rac := ea.RunnerActionConfig()
						ilv.Id = &rac
					})
					if err != nil {
						panic(err)
					}
					av.Agent = ea.Agent()
					av.Object = oa.Object()
					av.StartTime = time.Now()
					av.TargetCollection = tl
				})
				rs.messageQueue <- jsonld.Marshal(ca)

				act = act.WithResult(result.Create("201", "RAN HOOK"))
				rs.messageQueue <- jsonld.Marshal(act)
				w.WriteHeader(202)
			default:
				panic("Not a runner action")
			}
		default:
			panic("ONLY EXECUTE")
		}
	default:
		http.Error(w, "Unsupported action", 400)
		return
	}
}
