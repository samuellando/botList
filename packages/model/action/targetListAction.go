package action

import (
	"fedilist/packages/jsonld"
	"fedilist/packages/model/list"
	"fedilist/packages/model/result"
	"fmt"
	"time"
)

type targetListAction struct {
	action           action
	targetCollection list.ItemList
}

type targetListActionValues struct {
	Agent            Agent
	Object           list.ItemList
	StartTime        time.Time
	EndTime          *time.Time
	Result           *result.Result
	TargetCollection list.ItemList
}

func createTargetListAction(fs ...func(*targetListActionValues)) targetListAction {
	v := targetListActionValues{}
	for _, f := range fs {
		f(&v)
	}
	action := createAction(func(av *actionValues) {
		av.Agent = v.Agent
		av.Object = v.Object
		av.StartTime = v.StartTime
		av.EndTime = v.EndTime
		av.Result = v.Result
	})
	return targetListAction{
		action:           action,
		targetCollection: v.TargetCollection,
	}
}

type marshaledTargetListAction struct {
	Type             string         `json:"@type"`
	Agent            Agent         `json:"http://schema.org/agent"`
	Object           list.ItemList  `json:"http://schema.org/object"`
	StartTime        time.Time      `json:"http://schema.org/startTime"`
	EndTime          *time.Time     `json:"http://schema.org/endTime,omitempty"`
	Result           *result.Result `json:"http://schema.org/result,omitempty"`
	TargetCollection list.ItemList  `json:"http://schema.org/targetCollection"`
	Signature        string         `json:"http://fedilist.com/signature,omitempty"`
}

func (a targetListAction) marshal() marshaledTargetListAction {
	return marshaledTargetListAction{
		Agent:            a.action.agent,
		Object:           a.action.object,
		StartTime:        a.action.startTime,
		EndTime:          a.action.endTime,
		Result:           a.action.result,
		TargetCollection: a.targetCollection,
		Signature:        a.action.signature,
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
		t, err := list.Parse(json)
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
