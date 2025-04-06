package person

import (
	"encoding/json"
	"fedilist/packages/parser/jsonld"
	"fedilist/packages/parser/list"
	"fmt"
)

type Person struct {
	Id          string
	Name        *string
	Description *string
	List        []list.ItemList
}

func (p Person) MarshalJSON() ([]byte, error) {
	type External struct {
		Type        string          `json:"@type"`
		Id          string          `json:"@id"`
		Name        *string         `json:"http://schema.org/Name,omitempty"`
		Description *string         `json:"http://schema.org/Description,omitempty"`
		List        []list.ItemList `json:"http://fedilist.com/List,omitempty,omitempty"`
	}
	return json.Marshal(External{
		Type:        "http://schema.org/Person",
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		List:        p.List,
	})
}

var DB = make(map[string]Person)

func CreatePerson(domain, name, description string) Person {
	p := Person{
		Id:          domain + "/user/" + name,
		Name:        &name,
		Description: &description,
		List:        make([]list.ItemList, 0),
	}
	DB[p.Id] = p
	return p
}

func (p *Person) AddList(l list.ItemList) {
	p.List = append(p.List, l)
	DB[p.Id] = *p
}

func LoadPerson(json map[string]any) (Person, error) {
	if jsonld.GetType(json) != "http://schema.org/Person" {
		return Person{}, fmt.Errorf("Cannot load non person")
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)

	id := jsonld.GetId(json)
	if id == nil {
		return Person{}, fmt.Errorf("An id must be provided for person")
	}
	var name *string
	if v, ok := strs["name"]; ok {
		name = &v
	}
	var description *string
	if v, ok := strs["description"]; ok {
		description = &v
	}

	return Person{
		Id:          *id,
		Name:        name,
		Description: description,
	}, nil
}
