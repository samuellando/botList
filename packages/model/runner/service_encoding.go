package runner

import (
	"encoding/json"
	"fedilist/packages/jsonld"
	"fmt"
)

func (t Service) MarshalJSON() ([]byte, error) {
	type External struct {
		Type  string `json:"@type"`
		Name  string `json:"http://schema.org/name"`
		Schema string `json:"http://fedilist.com/schema"`

	}
	return json.Marshal(External{
		Type:               "http://fedilist.com/Service",
        Name: t.Name,
        Schema: t.Schema,
	})
}

func ParseService(json map[string]any) (Service, error) {
	if jsonld.GetType(json) != "http://fedilist.com/Service" {
		return Service{}, fmt.Errorf("Type must be Service")
	}
	fediOrgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)

	var name string
	if v, ok := strs["name"]; ok {
		name = v
	}

	strs = jsonld.GetBaseTypeValues[string](fediOrgValues)
	var schema string
	if v, ok := strs["schema"]; ok {
		schema = v
	}

	return Service{
        Name: name,
        Schema: schema,
	}, nil
}
