package list

import (
	"encoding/json"
	"fedilist/packages/parser/jsonld"
	"fmt"
)

func (t Runner) MarshalJSON() ([]byte, error) {
	type External struct {
		Type  string `json:"@type"`
		Id    string `json:"@id"`
		Name  string `json:"http://schema.org/name"`
		Inbox string `json:"https://fedilist.com/inbox"`
	}
	return json.Marshal(External{
		Type:               "http://fedilist.com/ActionHook",
        Id: t.id,
        Name: t.name,
        Inbox: t.inbox,
	})
}

func ParseRunner(json map[string]any) (Runner, error) {
	hook, err := parseHook(json)

	orgValues := jsonld.GetNamespaceValues(json, "https://fedilist.com")
	if err != nil {
		return Runner{}, err
	}
	strs := jsonld.GetBaseTypeValues[string](orgValues)
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
