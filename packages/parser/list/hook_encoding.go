package list

import (
	"encoding/json"
	"fedilist/packages/parser/jsonld"
	"fmt"
)

func (t ActionHook) MarshalJSON() ([]byte, error) {
	type External struct {
		Type               string   `json:"@type"`
		OnActionType       []string `json:"https://fedilist.com/onActionType"`
		RunnerAction       string   `json:"https://fedilist.com/runnerAction"`
		RunnerActionConfig string   `json:"https://fedilist.com/runnerActionConfig"`
	}
	return json.Marshal(External{
		Type:               "http://fedilist.com/ActionHook",
		OnActionType:       t.onActionType,
		RunnerAction:       t.hook.runnerAction,
		RunnerActionConfig: t.hook.runnerActionConfig,
	})
}

func (t CronHook) MarshalJSON() ([]byte, error) {
	type External struct {
		Type               string `json:"@type"`
		CronTab            string `json:"https://fedilist.com/cronTab"`
		RunnerAction       string `json:"https://fedilist.com/runnerAction"`
		RunnerActionConfig string `json:"https://fedilist.com/runnerActionConfig"`
	}
	return json.Marshal(External{
		Type:               "http://fedilist.com/CronHook",
		CronTab:             t.cronTab,
		RunnerAction:       t.hook.runnerAction,
		RunnerActionConfig: t.hook.runnerActionConfig,
	})
}

func ParseHook(json map[string]any) (Hook, error) {
    switch jsonld.GetType(json) {
    case "https://fedilist.com/ActionHook":
        return parseActionHook(json)
    case "https://fedilist.com/CronHook":
        return parseCronHook(json)
    default:
		return CronHook{}, fmt.Errorf("Type must be known hook type")
    }
}

func parseHook(json map[string]any) (hook, error) {
	orgValues := jsonld.GetNamespaceValues(json, "https://fedilist.com")
	strs := jsonld.GetBaseTypeValues[string](orgValues)
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
    return hook{
        runnerAction: runnerAction,
        runnerActionConfig: runnerActionConfig,
    }, nil
}

func parseCronHook(json map[string]any) (CronHook, error) {
    hook, err := parseHook(json)

	orgValues := jsonld.GetNamespaceValues(json, "https://fedilist.com")
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
        hook: hook,
        cronTab: cronTab,
	}, nil
}

func parseActionHook(json map[string]any) (ActionHook, error) {
    hook, err := parseHook(json)

	orgValues := jsonld.GetNamespaceValues(json, "https://fedilist.com")
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
        hook: hook,
        onActionType: onActionType,
	}, nil
}
