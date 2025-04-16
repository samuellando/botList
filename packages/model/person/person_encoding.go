package person

import (
	"encoding/json"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/list"
	"fmt"
)

func (p Person) MarshalJSON() ([]byte, error) {
	type External struct {
		Type        string          `json:"@type"`
		Id          string          `json:"@id,omitempty"`
		Name        string          `json:"http://schema.org/name,omitempty"`
		Description string          `json:"http://schema.org/description,omitempty"`
		Key         string          `json:"http://fedilist.com/key,omitempty"`
		List        []list.ItemList `json:"http://fedilist.com/list,omitempty"`
	}
	return json.Marshal(External{
		Type:        "http://schema.org/Person",
		Id:          p.Id(),
		Name:        p.Name(),
		Description: p.Description(),
		List:        p.List(),
		Key:         p.Key(),
	})
}

func LoadPerson(json map[string]any) (Person, error) {
	if jsonld.GetType(json) != "http://schema.org/Person" {
		return Person{}, fmt.Errorf("Cannot load non person")
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)

	var id string
	pid := jsonld.GetId(json)
	if pid != nil {
		id = *pid
	}
	var name string
	if v, ok := strs["name"]; ok {
		name = v
	}
	var description string
	if v, ok := strs["description"]; ok {
		description = v
	}

	return Person{
		id:          id,
		name:        name,
		description: description,
	}, nil
}
