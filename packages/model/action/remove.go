package action

import (
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"time"
)

type Remove struct {
	targetListAction targetListAction
}

func (a Remove) Signature() string {
	return a.targetListAction.action.signature
}

func (a Remove) Agent() Agent {
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

func (a Remove) TargetId() *string {
	id := a.targetListAction.targetCollection.Id()
	return &id
}

func (a Remove) Sign(s string) Action {
	a.targetListAction.action.signature = s
	return a
}

func (a Remove) TargetCollection() list.ItemList {
	return a.targetListAction.targetCollection
}

func (a Remove) WithResult(r result.Result) Action {
	t := time.Now()
	a.targetListAction.action.result = &r
	a.targetListAction.action.endTime = &t
	return a
}
