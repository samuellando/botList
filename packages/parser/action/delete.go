package action

import (
	"fedilist/packages/parser/jsonld"
	"fedilist/packages/parser/person"
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/result"
	"fmt"
    "time"
)

type Delete struct {
    targetListAction targetListAction
}

func (a Delete) Agent() person.Person {
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

func (a Delete) Serialize() []byte {
	return []byte{}
}

func (a Delete) TargetId() *string {
	return a.targetListAction.targetCollection.Id()
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

