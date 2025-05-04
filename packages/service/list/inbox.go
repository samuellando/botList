package list

import (
	"fedilist/packages/model/list"
	"fedilist/packages/model/action"
	"fedilist/packages/util"
	"net/http"
	"fmt"
)

func (ls ListService) handleInbox(list list.ItemList, w http.ResponseWriter, req *http.Request) {
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
		case action.Prepend:
			ls.Prepend(w, act)
		case action.Remove:
			ls.Remove(w, act)
		case action.Update:
			ls.Update(w, act)
		default:
			http.Error(w, "Unsupported action", 400)
			return
		}
}
