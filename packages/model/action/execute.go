package action

import (
	"encoding/json"
	"fedilist/packages/model/result"
	"fedilist/packages/model/runner"
	"time"
)

type Execute struct {
	agent              Agent
	object             Action
	signature          string
	targetRunner       runner.Runner
	runnerAction       string
	runnerActionConfig string
	startTime          time.Time
	endTime            *time.Time
	result             *result.Result
}

type ExecuteValues struct {
	Agent              Agent
	Object             Action
	Signature          string
	TargetRunner       runner.Runner
	RunnerAction       string
	RunnerActionConfig string
	StartTime          time.Time
	EndTime            *time.Time
	Result             *result.Result
}

func CreateExecute(fs ...func(*ExecuteValues)) Execute {
	v := ExecuteValues{}
	for _, f := range fs {
		f(&v)
	}
	return Execute{
		agent:              v.Agent,
		object:             v.Object,
		signature:          v.Signature,
		targetRunner:       v.TargetRunner,
		runnerAction:       v.RunnerAction,
		runnerActionConfig: v.RunnerActionConfig,
		startTime:          v.StartTime,
		endTime:            v.EndTime,
		result:             v.Result,
	}
}

func (a Execute) MarshalJSON() ([]byte, error) {
	type External struct {
		Type               string         `json:"@type"`
		Agent              Agent          `json:"http://schema.org/agent"`
		Object             Action         `json:"http://schema.org/object"`
		Signature          string         `json:"http://fedilist.com/signature"`
		TargetRunner       runner.Runner  `json:"http://fedilist.com/targetRunner"`
		RunnerAction       string         `json:"http://fedilist.com/runnerAction"`
		RunnerActionConfig string         `json:"http://fedilist.com/runnerActionConfig"`
		StartTime          time.Time      `json:"http://schema.org/startTime"`
		EndTime            *time.Time     `json:"http://schema.org/endTime,omitempty"`
		Result             *result.Result `json:"http://schema.org/result,omitempty"`
	}
	return json.Marshal(External{
		Type:               "http://fedilist.com/ExecuteAction",
		Agent:              a.Agent(),
		Object:             a.Object(),
		Signature:          a.Signature(),
		TargetRunner:       a.TargetRunner(),
		RunnerAction:       a.RunnerAction(),
		RunnerActionConfig: a.RunnerActionConfig(),
		StartTime:          a.StartTime(),
		EndTime:            a.EndTime(),
		Result:             a.Result(),
	})
}

func (a Execute) Agent() Agent {
	return a.agent
}

func (a Execute) Object() Action {
	return a.object
}

func (a Execute) Signature() string {
	return a.signature
}

func (a Execute) Sign(s string) Action {
	a.signature = s
	return a
}

func (a Execute) StartTime() time.Time {
	return a.startTime
}

func (a Execute) EndTime() *time.Time {
	return a.endTime
}

func (a Execute) Result() *result.Result {
	return a.result
}

func (a Execute) TargetId() *string {
	id := a.targetRunner.Id()
	return &id
}

func (a Execute) TargetRunner() runner.Runner {
	return a.targetRunner
}

func (a Execute) WithResult(r result.Result) Action {
	t := time.Now()
	a.result = &r
	a.endTime = &t
	return a
}

func (a Execute) RunnerAction() string {
	return a.runnerAction
}

func (a Execute) RunnerActionConfig() string {
	return a.runnerActionConfig
}
