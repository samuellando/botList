package action

import (
	"fedilist/packages/parser/jsonld"
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/person"
	"fedilist/packages/parser/result"
	"fmt"
	"time"
)

type Remove struct {
	targetListAction targetListAction
	atIndex          int
}

func (a Remove) Agent() person.Person {
	return a.targetListAction.action.agent
}

func (a Remove) Object() list.ItemList {
	return a.targetListAction.action.object
}

func (a Remove) StartTime() time.Time {
	return a.targetListAction.action.startTime
}

func (a Remove) EndTime() *time.Time {
	return a.targetListAction.action.endTime
}

func (a Remove) Result() *result.Result {
	return a.targetListAction.action.result
}

func (a Remove) AtIndex() int {
	return a.atIndex
}

func (a Remove) Serialize() []byte {
	return []byte{}
}

func (a Remove) TargetId() *string {
	return a.targetListAction.targetCollection.Id()
}


func parseRemove(json map[string]any) (Remove, error) {
	if jsonld.GetType(json) != "http://schema.org/RemoveAction" {
		return Remove{}, fmt.Errorf("Wrong @type")
	}
	tla, err := parseTargetListAction(json)
	if err != nil {
		return Remove{}, err
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "https://fedilist.com")
	ints := jsonld.GetBaseTypeValues[float64](schemaOrgValues)
	var atIndex int
	if i, ok := ints["atIndex"]; ok {
		atIndex = int(i)
	}
	return Remove{targetListAction: tla, atIndex: atIndex}, nil
}
