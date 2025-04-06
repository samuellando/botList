package action

import (
	"encoding/json"
	"fedilist/packages/parser/jsonld"
	"fedilist/packages/parser/list"
	"fedilist/packages/parser/person"
	"fedilist/packages/parser/result"
	"fmt"
	"time"
)

type action struct {
    agent     person.Person
	object    list.ItemList
	startTime time.Time
	endTime   *time.Time
	result    *result.Result
}

func (a action) MarshalJSON() ([]byte, error) {
	type External struct {
        Agent     person.Person   `json:"http://schema.org/Agent"`
		Object    list.ItemList   `json:"http://schema.org/Object"`
		StartTime time.Time       `json:"http://schema.org/StartTime"`
		EndTime   *time.Time      `json:"http://schema.org/EndTime,omitempty"`
		Result    *result.Result  `json:"http://schema.org/Result,omitempty"`
	}
	return json.Marshal(External{
		Agent:     a.agent,
		Object:    a.object,
		StartTime: a.startTime,
		EndTime:   a.endTime,
		Result:    a.result,
	})
}

type Action interface {
	Agent() person.Person
    Result() *result.Result
	Serialize() ([]byte, error)
    TargetId() *string
}

func parseAction(json map[string]any) (action, error) {
	var err error
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	objs := jsonld.GetCompositeTypeValues(schemaOrgValues)

	var agent person.Person
	if json, ok := objs["agent"]; ok {
		agent, err = person.LoadPerson(json)
		if err != nil {
			return action{}, err
		}
	} else {
		return action{}, fmt.Errorf("Actions must have an agent")
	}

	var object list.ItemList
	if json, ok := objs["object"]; ok {
		object, err = list.LoadItemList(json)
		if err != nil {
			return action{}, err
		}
	} else {
		return action{}, fmt.Errorf("Actions must have an object")
	}

	times := jsonld.GetBaseTypeValues[string](schemaOrgValues)
	var startTime time.Time
	if st, ok := times["startTime"]; ok {
		t, err := time.Parse(time.RFC3339, st)
		if err != nil {
			return action{}, fmt.Errorf("Invalid start time format")
		}
		startTime = t
	} else {
		return action{}, fmt.Errorf("Actions must have a start time")
	}

	var endTime *time.Time
	if et, ok := times["endTime"]; ok {
		t, err := time.Parse(time.RFC3339, et)
		if err != nil {
			return action{}, fmt.Errorf("Invalid end time format")
		}
		endTime = &t
	}

	var res *result.Result
	if json, ok := objs["result"]; ok {
		r, err := result.LoadResult(json)
		if err != nil {
			return action{}, err
		}
		res = &r
	}

	return action{
		agent:     agent,
		object:    object,
		startTime: startTime,
		endTime:   endTime,
		result:    res,
	}, nil
}
