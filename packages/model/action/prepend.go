package action

import (
	"fedilist/packages/jsonld"
	"fedilist/packages/model/list"
	"fedilist/packages/model/person"
	"fedilist/packages/model/result"
	"fmt"
	"time"
)

type Prepend struct {
	targetListAction targetListAction
}

func (a Prepend) Agent() person.Person {
	return a.targetListAction.action.agent
}

func (a Prepend) Object() list.ItemList {
	return a.targetListAction.action.object
}

func (a Prepend) StartTime() time.Time {
	return a.targetListAction.action.startTime
}

func (a Prepend) EndTime() *time.Time {
	return a.targetListAction.action.endTime
}

func (a Prepend) Result() *result.Result {
	return a.targetListAction.action.result
}

func (a Prepend) TargetId() *string {
	return a.targetListAction.targetCollection.Id()
}

func (a Prepend) WithResult(r result.Result) Action {
	t := time.Now()
	a.targetListAction.action.result = &r
	a.targetListAction.action.endTime = &t
	return a
}


func parsePrepend(json map[string]any) (Prepend, error) {
	if jsonld.GetType(json) != "http://schema.org/PrependAction" {
		return Prepend{}, fmt.Errorf("Wrong @type")
	}
	tla, err := parseTargetListAction(json)
	if err != nil {
		return Prepend{}, err
	}
	return Prepend{targetListAction: tla}, nil
}
