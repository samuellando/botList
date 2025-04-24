package action

import (
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"time"
)

type Create struct {
	action action
}

func (a Create) Signature() string {
	return a.action.signature
}

func (a Create) Agent() Agent {
	return a.action.agent
}

func (a Create) Object() list.ItemList {
	return a.action.object
}

func (a Create) StartTime() time.Time {
	return a.action.startTime
}

func (a Create) EndTime() *time.Time {
	return a.action.endTime
}

func (a Create) Result() *result.Result {
	return a.action.result
}

func (a Create) TargetId() *string {
	return nil
}

func (a Create) Sign(s string) Action {
	a.action.signature = s
	return a
}

func (a Create) WithResult(r result.Result) Action {
	t := time.Now()
	a.action.result = &r
	a.action.endTime = &t
	return a
}
