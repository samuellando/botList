package action

import (
	"fedilist/packages/parser/jsonld"
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/person"
	"fedilist/packages/parser/result"
	"fmt"
	"time"
)

type targetListAction struct {
	action           action
	targetCollection list.ItemList
}

type marshaledTargetListAction struct {
	Type             string         `json:"@type"`
	Agent            person.Person  `json:"http://schema.org/agent"`
	Object           list.ItemList  `json:"http://schema.org/object"`
	StartTime        time.Time      `json:"http://schema.org/startTime"`
	EndTime          *time.Time     `json:"http://schema.org/endTime,omitempty"`
	Result           *result.Result `json:"http://schema.org/result,omitempty"`
	TargetCollection list.ItemList  `json:"http://schema.org/targetCollection"`
}

func (a targetListAction) marshal() marshaledTargetListAction {
	return marshaledTargetListAction {
		Agent:            a.action.agent,
		Object:           a.action.object,
		StartTime:        a.action.startTime,
		EndTime:          a.action.endTime,
		Result:           a.action.result,
		TargetCollection: a.targetCollection,
	}
}

func parseTargetListAction(json map[string]any) (targetListAction, error) {
	action, err := parseAction(json)
	if err != nil {
		return targetListAction{}, err
	}

	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	objs := jsonld.GetCompositeTypeValues(schemaOrgValues)

	var targetCollection list.ItemList
	if json, ok := objs["targetCollection"]; ok {
		t, err := list.LoadItemList(json)
		if err != nil {
			return targetListAction{}, err
		}
		targetCollection = t
	} else {
		return targetListAction{}, fmt.Errorf("targetListActions must have target collection")
	}
	return targetListAction{
		action:           action,
		targetCollection: targetCollection,
	}, nil
}
