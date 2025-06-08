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
		Inbox       string          `json:"http://fedilist.com/inbox,omitempty"`
		Outbox      string          `json:"http://fedilist.com/outbox,omitempty"`
	}
	return json.Marshal(External{
		Type:        "http://schema.org/Person",
		Id:          p.Id(),
		Name:        p.Name(),
		Description: p.Description(),
		List:        p.List(),
		Key:         p.Key(),
		Inbox:       p.Inbox(),
		Outbox:      p.Outbox(),
	})
}

func (p *Person) UnmarshalJSON(data []byte) error {
	type External struct {
		Type        string          `json:"@type"`
		Id          string          `json:"@id,omitempty"`
		Name        string          `json:"http://schema.org/name,omitempty"`
		Description string          `json:"http://schema.org/description,omitempty"`
		Key         string          `json:"http://fedilist.com/key,omitempty"`
		List        []list.ItemList `json:"http://fedilist.com/list,omitempty"`
		Inbox       string          `json:"http://fedilist.com/inbox,omitempty"`
		Outbox      string          `json:"http://fedilist.com/outbox,omitempty"`
	}
	var ext External
	if err := json.Unmarshal(data, &ext); err != nil {
		return err
	}
	if ext.Type != "http://schema.org/Person" {
		return fmt.Errorf("invalid type: %s", ext.Type)
	}
	p.id = ext.Id
	p.name = ext.Name
	p.description = ext.Description
	p.key = ext.Key
	return nil
}

func LoadPerson(json map[string]any) (Person, error) {
	if jsonld.GetType(json) != "http://schema.org/Person" {
		return Person{}, fmt.Errorf("Cannot load non person")
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)

	id := jsonld.GetId(json)
	var name string
	if v, ok := strs["name"]; ok {
		name = v
	}
	var description string
	if v, ok := strs["description"]; ok {
		description = v
	}

	fediOrgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	strs = jsonld.GetBaseTypeValues[string](fediOrgValues)

	var key string
	if v, ok := strs["key"]; ok {
		key = v
	}

	return Person{
		id:          id,
		name:        name,
		description: description,
		key: key,
	}, nil
}
