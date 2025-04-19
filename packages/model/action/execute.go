package action

import (
	"encoding/json"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/person"
	"fedilist/packages/model/result"
	"fedilist/packages/model/runner"
	"fmt"
	"time"
)

type Execute struct {
	agent              person.Person
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
	Agent              person.Person
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
		Agent              person.Person  `json:"http://schema.org/agent"`
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

func (a Execute) Agent() person.Person {
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
	return a.targetRunner.Id()
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

func parseExecute(json map[string]any) (Execute, error) {
	if jsonld.GetType(json) != "http://fedilist.com/ExecuteAction" {
		return Execute{}, fmt.Errorf("Wrong @type")
	}
	var err error
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	fediOrgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	objs := jsonld.GetCompositeTypeValues(schemaOrgValues)

	var agent person.Person
	if json, ok := objs["agent"]; ok {
		agent, err = person.LoadPerson(json)
		if err != nil {
			return Execute{}, err
		}
	} else {
		return Execute{}, fmt.Errorf("Actions must have an agent")
	}

	var object Action
	if json, ok := objs["object"]; ok {
		object, err = Parse(json)
		if err != nil {
			return Execute{}, err
		}
	}

	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)
	var startTime time.Time
	if st, ok := strs["startTime"]; ok {
		t, err := time.Parse(time.RFC3339, st)
		if err != nil {
			return Execute{}, fmt.Errorf("Invalid start time format")
		}
		startTime = t
	} else {
		return Execute{}, fmt.Errorf("Actions must have a start time")
	}

	var endTime *time.Time
	if et, ok := strs["endTime"]; ok {
		t, err := time.Parse(time.RFC3339, et)
		if err != nil {
			return Execute{}, fmt.Errorf("Invalid end time format")
		}
		endTime = &t
	}

	strs = jsonld.GetBaseTypeValues[string](fediOrgValues)

	var runnerAction string
	if v, ok := strs["runnerAction"]; ok {
		runnerAction = v
	} else {
		return Execute{}, fmt.Errorf("Must have a runner action")
	}

	var runnerActionConfig string
	if v, ok := strs["runnerActionConfig"]; ok {
		runnerActionConfig = v
	} else {
		return Execute{}, fmt.Errorf("Must have a runner action config")
	}

	var res *result.Result
	if json, ok := objs["result"]; ok {
		r, err := result.LoadResult(json)
		if err != nil {
			return Execute{}, err
		}
		res = &r
	}

	objs = jsonld.GetCompositeTypeValues(fediOrgValues)

	var r runner.Runner
	if json, ok := objs["targetRunner"]; ok {
		r, err = runner.Parse(json)
		if err != nil {
			return Execute{}, err
		}
	} else {
		return Execute{}, fmt.Errorf("Execute must have targetRunners")
	}

	return Execute{
		agent:              agent,
		object:             object,
		startTime:          startTime,
		endTime:            endTime,
		result:             res,
		targetRunner:       r,
		runnerAction:       runnerAction,
		runnerActionConfig: runnerActionConfig,
	}, nil
}
