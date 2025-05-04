package action

import (
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"time"
)

type Update struct {
	targetListAction targetListAction
}

func (a Update) Signature() string {
	return a.targetListAction.action.signature
}

func (a Update) Agent() Agent {
	return a.targetListAction.action.agent
}

func (a Update) Object() list.ItemList {
	return a.targetListAction.action.object
}

func (a Update) StartTime() time.Time {
	return a.targetListAction.action.startTime
}

func (a Update) EndTime() *time.Time {
	return a.targetListAction.action.endTime
}

func (a Update) Result() *result.Result {
	return a.targetListAction.action.result
}

func (a Update) TargetId() *string {
	id := a.targetListAction.targetCollection.Id()
	return &id
}

func (a Update) TargetCollection() list.ItemList {
	return a.targetListAction.targetCollection
}

func (a Update) Sign(s string) Action {
	a.targetListAction.action.signature = s
	return a
}

func (a Update) WithResult(r result.Result) Action {
	t := time.Now()
	a.targetListAction.action.result = &r
	a.targetListAction.action.endTime = &t
	return a
}
