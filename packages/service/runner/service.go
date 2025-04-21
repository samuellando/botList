package runner

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/action"
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"fedilist/packages/model/runner"
	"fedilist/packages/util"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	messageQueue chan []byte
	runner       runner.Runner
	key          []byte
}

func Create(id string, q chan []byte) Service {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
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
		rv.Key = base64.StdEncoding.EncodeToString(publicKey)
	})
	if err != nil {
		panic(err)
	}
	return Service{
		messageQueue: q,
		runner:       r,
		key:          privateKey.Seed(),
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
		json, err := util.GetBodyJsonld(req.Body)
		if err != nil {
			panic(err)
		}
		act, err := action.Parse(json)
		if err != nil {
			panic(err)
		}
		var resp *http.Response
		if act.Result() == nil {
			resp, err = http.Get(act.Agent().Id())
		} else {
			resp, err = http.Get(*act.TargetId())
		}
		personJson, err := util.GetBodyJsonld(resp.Body)
		if err != nil {
			panic(err)
		}
		agent, err := action.ParseAgent(personJson)
		if err != nil {
			panic(err)
		}
		valid, err := util.VerifySignature(act, agent.Key())
		if err != nil {
			panic(err)
		}
		if !valid {
			http.Error(w, "Could not verify message signature", 403)
			return
		}
		if act.Result() != nil {
			fmt.Println("Runner GOT RESULT")
			w.WriteHeader(202)
			return
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
						ilv.Id = rac
					})
					if err != nil {
						panic(err)
					}
					av.Agent = rs.runner
					av.Object = oa.Object()
					av.StartTime = time.Now()
					av.TargetCollection = tl
				})
				rs.messageQueue <- jsonld.Marshal(util.Sign(ca, rs.key))

				act = act.WithResult(result.Create("201", "RAN HOOK"))
				rs.messageQueue <- jsonld.Marshal(util.Sign(act, rs.key))
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
