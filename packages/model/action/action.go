package action

import (
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"time"
)

type action struct {
	agent     Agent
	object    list.ItemList
	startTime time.Time
	endTime   *time.Time
	result    *result.Result
	signature string
}

type Action interface {
	Agent() Agent
	Result() *result.Result
	TargetId() *string
	WithResult(result.Result) Action
	Signature() string
	Sign(string) Action
}

type Agent interface {
	Id() string
	Key() string
}

type actionValues struct {
	Agent     Agent
	Object    list.ItemList
	StartTime time.Time
	EndTime   *time.Time
	Result    *result.Result
	Signature string
}

func createAction(fs ...func(*actionValues)) action {
	v := actionValues{}
	for _, f := range fs {
		f(&v)
	}
	return action{
		agent:     v.Agent,
		object:    v.Object,
		startTime: v.StartTime,
		endTime:   v.EndTime,
		result:    v.Result,
		signature: v.Signature,
	}
}
