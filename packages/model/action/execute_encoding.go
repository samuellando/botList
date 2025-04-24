package action

import (
	"fedilist/packages/jsonld"
	"fedilist/packages/model/result"
	"fedilist/packages/model/runner"
	"fmt"
	"time"
)


func parseExecute(json map[string]any) (Execute, error) {
	if jsonld.GetType(json) != "http://fedilist.com/ExecuteAction" {
		return Execute{}, fmt.Errorf("Wrong @type")
	}
	var err error
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	fediOrgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	objs := jsonld.GetCompositeTypeValues(schemaOrgValues)

	var agent Agent
	if json, ok := objs["agent"]; ok {
		agent, err = ParseAgent(json)
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
	
	var signature string
	if v, ok := strs["signature"]; ok {
		signature = v
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
		signature: signature,
	}, nil
}
