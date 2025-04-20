package action

import (
	"fedilist/packages/jsonld"
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"fmt"
    "time"
)

type Delete struct {
    targetListAction targetListAction
}

func (a Delete) Signature() string {
	return a.targetListAction.action.signature
}

func (a Delete) Agent() Agent {
	return a.targetListAction.action.agent
}

func (a Delete) Object() list.ItemList {
	return a.targetListAction.action.object
}

func (a Delete) StartTime() time.Time {
	return a.targetListAction.action.startTime
}

func (a Delete) EndTime() *time.Time {
	return a.targetListAction.action.endTime
}

func (a Delete) Result() *result.Result{
	return a.targetListAction.action.result
}

func (a Delete) TargetId() *string {
	id := a.targetListAction.targetCollection.Id()
	return &id
}

func (a Delete) Sign(s string) Action {
	a.targetListAction.action.signature = s
	return a
}

func (a Delete) WithResult(r result.Result) Action {
	t := time.Now()
	a.targetListAction.action.result = &r
	a.targetListAction.action.endTime = &t
	return a
}

func parseDelete(json map[string]any) (Delete, error) {
    if jsonld.GetType(json) != "http://schema.org/DeleteAction" {
        return Delete{}, fmt.Errorf("Wrong @type")
    }
    tla, err := parseTargetListAction(json)
	if err != nil {
		return Delete{}, err
	}
    return Delete{targetListAction: tla}, nil
}

