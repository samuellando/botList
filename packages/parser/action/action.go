package action

import (
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/person"
	"fedilist/packages/parser/result"
	"time"
)

type action struct {
    agent     person.Person
	object    list.ItemList
	startTime time.Time
	endTime   *time.Time
	result    *result.Result
}

type Action interface {
	Agent() person.Person
    Result() *result.Result
    TargetId() *string
    WithResult(result.Result) Action
}

type actionValues struct {
	Agent              person.Person
	Object             list.ItemList
	StartTime          time.Time
	EndTime            *time.Time
	Result             *result.Result
}

func createAction(fs ...func(*actionValues)) action {
	v := actionValues{}
	for _, f := range fs {
		f(&v)
	}
	return action{
		agent:              v.Agent,
		object:             v.Object,
		startTime:          v.StartTime,
		endTime:            v.EndTime,
		result:             v.Result,
	}
}

