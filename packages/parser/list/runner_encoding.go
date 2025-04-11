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
        Service []Service `json:"http://fedilist.com/service"`

	}
	return json.Marshal(External{
		Type:               "http://fedilist.com/Runner",
        Id: t.id,
        Name: t.name,
        Inbox: t.inbox,
        Service: t.service,
	})
}

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

    arrs := jsonld.GetCompositeTypeArrayValues(schemaOrgValues)
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

	return Runner{
        id: *id,
        name: name,
        inbox: inbox,
	}, nil
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
