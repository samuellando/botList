package action

import (
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"time"
)

type Append struct {
	targetListAction targetListAction
}

type AppendValues struct {
	Agent            Agent
	Object           list.ItemList
	StartTime        time.Time
	EndTime          *time.Time
	Result           *result.Result
	TargetCollection list.ItemList
}

func CreateAppend(fs ...func(*AppendValues)) Append {
	v := AppendValues{}
	for _, f := range fs {
		f(&v)
	}
    targetListAction := createTargetListAction(func(tlav *targetListActionValues) {
		tlav.Agent = v.Agent
		tlav.Object = v.Object
		tlav.StartTime = v.StartTime
		tlav.EndTime = v.EndTime
		tlav.Result = v.Result
        tlav.TargetCollection = v.TargetCollection
    })
	return Append{
        targetListAction: targetListAction,
	}
}

func (a Append) Agent() Agent {
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

func (a Append) WithResult(r result.Result) Action {
	t := time.Now()
	a.targetListAction.action.result = &r
	a.targetListAction.action.endTime = &t
	return a
}

func (a Append) Signature() string {
	return a.targetListAction.action.signature
}

func (a Append) Sign(s string) Action {
	a.targetListAction.action.signature = s
	return a
}

func (a Append) TargetId() *string {
	id := a.targetListAction.targetCollection.Id()
	return &id
}

func (a Append) TargetCollection() list.ItemList {
	return a.targetListAction.targetCollection
}
