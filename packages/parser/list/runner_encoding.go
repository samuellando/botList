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
		Inbox string `json:"http://fedilist.com/inbox"`
	}
	return json.Marshal(External{
		Type:               "http://fedilist.com/Runner",
        Id: t.id,
        Name: t.name,
        Inbox: t.inbox,
	})
}

func ParseRunner(json map[string]any) (Runner, error) {
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

	strs = jsonld.GetBaseTypeValues[string](fediOrgValues)
	var inbox string
	if v, ok := strs["inbox"]; ok {
		inbox = v
	}

	return Runner{
        id: *id,
        name: name,
        inbox: inbox,
	}, nil
}
