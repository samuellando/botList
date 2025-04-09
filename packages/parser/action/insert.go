package action

import (
	"fedilist/packages/parser/jsonld"
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/person"
	"fedilist/packages/parser/result"
	"fmt"
	"time"
)

type Insert struct {
	targetListAction targetListAction
	atIndex          int
}

func (a Insert) Agent() person.Person {
	return a.targetListAction.action.agent
}

func (a Insert) Object() list.ItemList {
	return a.targetListAction.action.object
}

func (a Insert) StartTime() time.Time {
	return a.targetListAction.action.startTime
}

func (a Insert) EndTime() *time.Time {
	return a.targetListAction.action.endTime
}

func (a Insert) Result() *result.Result {
	return a.targetListAction.action.result
}

func (a Insert) AtIndex() int {
	return a.atIndex
}

func (a Insert) TargetId() *string {
	return a.targetListAction.targetCollection.Id()
}

func (a Insert) WithResult(r result.Result) Action {
	t := time.Now()
	a.targetListAction.action.result = &r
	a.targetListAction.action.endTime = &t
	return a
}


func parseInsert(json map[string]any) (Insert, error) {
	if jsonld.GetType(json) != "http://schema.org/InsertAction" {
		return Insert{}, fmt.Errorf("Wrong @type")
	}
	tla, err := parseTargetListAction(json)
	if err != nil {
		return Insert{}, err
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	ints := jsonld.GetBaseTypeValues[float64](schemaOrgValues)
	var atIndex int
	if i, ok := ints["atIndex"]; ok {
		atIndex = int(i)
	}
	return Insert{targetListAction: tla, atIndex: atIndex}, nil
}
