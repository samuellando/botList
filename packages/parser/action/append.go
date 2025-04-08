package action

import (
	"fedilist/packages/parser/jsonld"
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/person"
	"fedilist/packages/parser/result"
	"fmt"
	"time"
)

type Append struct {
	targetListAction targetListAction
}

func (a Append) Agent() person.Person {
	return a.targetListAction.action.agent
}

func (a Append) Object() list.ItemList {
	return a.targetListAction.action.object
}

func (a Append) StartTime() time.Time {
	return a.targetListAction.action.startTime
}

func (a Append) EndTime() *time.Time {
	return a.targetListAction.action.endTime
}

func (a Append) Result() *result.Result {
	return a.targetListAction.action.result
}

func (a *Append) AddResult(r result.Result) {
    t := time.Now()
	a.targetListAction.action.result = &r
	a.targetListAction.action.endTime = &t
}

func (a Append) TargetId() *string {
	return a.targetListAction.targetCollection.Id()
}

func (a Append) TargetCollection() list.ItemList {
	return a.targetListAction.targetCollection
}

func (a Append) Serialize() []byte {
    stla := a.targetListAction.marshal()
    stla.Type = "AppendAction"
	return jsonld.Marshal(CONTEXT, stla)
}

func parseAppend(json map[string]any) (Append, error) {
	if jsonld.GetType(json) != "http://schema.org/AppendAction" {
		return Append{}, fmt.Errorf("Wrong @type")
	}
	tla, err := parseTargetListAction(json)
	if err != nil {
		return Append{}, err
	}
	return Append{targetListAction: tla}, nil
}
