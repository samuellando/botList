package runner

import (
	"encoding/json"
	"fedilist/packages/jsonld"
	"fmt"
)

func (t Runner) MarshalJSON() ([]byte, error) {
	type External struct {
		Type    string    `json:"@type"`
		Id      string    `json:"@id"`
		Name    string    `json:"http://schema.org/name"`
		Inbox   string    `json:"http://fedilist.com/inbox"`
		Service []Service `json:"http://fedilist.com/service"`
		Key     string    `json:"http://fedilist.com/key"`
	}
	return json.Marshal(External{
		Type:    "http://fedilist.com/Runner",
		Id:      t.id,
		Name:    t.name,
		Inbox:   t.inbox,
		Service: t.service,
		Key:     t.Key(),
	})
}

func Parse(json map[string]any) (Runner, error) {
	if jsonld.GetType(json) != "http://fedilist.com/Runner" {
		return Runner{}, fmt.Errorf("Type must be Runner")
	}
	id := jsonld.GetId(json)
	fediOrgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)

	var name string
	if v, ok := strs["name"]; ok {
		name = v
	}

	arrs := jsonld.GetCompositeTypeArrayValues(fediOrgValues)
	services := make([]Service, 0)
	if a, ok := arrs["service"]; ok {
		for _, v := range a {
			s, err := ParseService(v)
			if err != nil {
				return Runner{}, err
			}
			services = append(services, s)
		}
	}

	strs = jsonld.GetBaseTypeValues[string](fediOrgValues)
	var inbox string
	if v, ok := strs["inbox"]; ok {
		inbox = v
	}
	var key string
	if v, ok := strs["key"]; ok {
		key = v
	}

	return Runner{
		id:    id,
		name:  name,
		inbox: inbox,
		key:   key,
		service: services,
	}, nil
}
