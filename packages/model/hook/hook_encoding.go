package hook

import (
	"encoding/json"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/runner"
	"fmt"
)

func (t ActionHook) MarshalJSON() ([]byte, error) {
	type External struct {
		Type               string   `json:"@type"`
		OnActionType       []string `json:"http://fedilist.com/onActionType"`
		Runner             runner.Runner   `json:"http://fedilist.com/runner"`
		RunnerAction       string   `json:"http://fedilist.com/runnerAction"`
		RunnerActionConfig string   `json:"http://fedilist.com/runnerActionConfig"`
	}
	return json.Marshal(External{
		Type:               "http://fedilist.com/ActionHook",
		OnActionType:       t.onActionType,
		Runner:             t.hook.runner,
		RunnerAction:       t.hook.runnerAction,
		RunnerActionConfig: t.hook.runnerActionConfig,
	})
}

func (t CronHook) MarshalJSON() ([]byte, error) {
	type External struct {
		Type               string `json:"@type"`
		Runner             runner.Runner `json:"http://fedilist.com/runner"`
		CronTab            string `json:"http://fedilist.com/cronTab"`
		RunnerAction       string `json:"http://fedilist.com/runnerAction"`
		RunnerActionConfig string `json:"http://fedilist.com/runnerActionConfig"`
	}
	return json.Marshal(External{
		Type:               "http://fedilist.com/CronHook",
		CronTab:            t.cronTab,
		Runner:             t.hook.runner,
		RunnerAction:       t.hook.runnerAction,
		RunnerActionConfig: t.hook.runnerActionConfig,
	})
}

func Parse(json map[string]any) (Hook, error) {
	switch jsonld.GetType(json) {
	case "http://fedilist.com/ActionHook":
		return parseActionHook(json)
	case "http://fedilist.com/CronHook":
		return parseCronHook(json)
	default:
		return CronHook{}, fmt.Errorf("Type must be known hook type")
	}
}

func parseHook(json map[string]any) (hook, error) {
	orgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	strs := jsonld.GetBaseTypeValues[string](orgValues)
	objs := jsonld.GetCompositeTypeValues(orgValues)
	var runnerAction string
	if v, ok := strs["runnerAction"]; ok {
		runnerAction = v
	} else {
		return hook{}, fmt.Errorf("Hooks must have a runner action")
	}
	var runnerActionConfig string
	if v, ok := strs["runnerActionConfig"]; ok {
		runnerActionConfig = v
	} else {
		return hook{}, fmt.Errorf("Hooks must have a runner action config")
	}
	var r runner.Runner
	var err error
	if o, ok := objs["runner"]; ok {
		r, err = runner.Parse(o)
		if err != nil {
			panic(err)
		}
	} else {
		return hook{}, fmt.Errorf("Hooks must have a runner")
	}
	return hook{
		runner: r,
		runnerAction:       runnerAction,
		runnerActionConfig: runnerActionConfig,
	}, nil
}

func parseCronHook(json map[string]any) (CronHook, error) {
	hook, err := parseHook(json)

	orgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	if err != nil {
		return CronHook{}, err
	}
	strs := jsonld.GetBaseTypeValues[string](orgValues)
	var cronTab string
	if v, ok := strs["cronTab"]; ok {
		cronTab = v
	} else {
		return CronHook{}, fmt.Errorf("Action hooks must have a onActionType")
	}

	return CronHook{
		hook:    hook,
		cronTab: cronTab,
	}, nil
}

func parseActionHook(json map[string]any) (ActionHook, error) {
	hook, err := parseHook(json)

	orgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	if err != nil {
		return ActionHook{}, err
	}
	strs := jsonld.GetBaseTypeArrayValues[string](orgValues)
	var onActionType []string
	if arr, ok := strs["onActionType"]; ok {
		onActionType = arr
	} else {
		return ActionHook{}, fmt.Errorf("Action hooks must have a onActionType")
	}

	return ActionHook{
		hook:         hook,
		onActionType: onActionType,
	}, nil
}
