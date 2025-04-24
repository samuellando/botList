package action

import (
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"time"
)

type Insert struct {
	targetListAction targetListAction
	atIndex          int
}

func (a Insert) Signature() string {
	return a.targetListAction.action.signature
}

func (a Insert) Agent() Agent {
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
	id := a.targetListAction.targetCollection.Id()
	return &id
}

func (a Insert) Sign(s string) Action {
	a.targetListAction.action.signature = s
	return a
}

func (a Insert) WithResult(r result.Result) Action {
	t := time.Now()
	a.targetListAction.action.result = &r
	a.targetListAction.action.endTime = &t
	return a
}
