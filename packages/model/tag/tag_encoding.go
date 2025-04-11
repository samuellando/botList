package tag

import (
	"encoding/json"
	"fedilist/packages/jsonld"
    "fmt"
)

func (t Tag) MarshalJSON() ([]byte, error) {
	type External struct {
		Type            string     `json:"@type"`
		Name            string    `json:"http://schema.org/dame"`
		Description     *string    `json:"http://schema.org/description,omitempty"`
	}
	return json.Marshal(External{
		Type:            "http://fedilist.com/Tag",
		Name:            t.name,
		Description:     t.description,
	})
}

func Parse(json map[string]any) (Tag, error) {
	if jsonld.GetType(json) != "http://fedilist.com/Tag" {
		return Tag{}, fmt.Errorf("Type must be Tag")
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)
	var name *string
	if v, ok := strs["name"]; ok {
		name = &v
	} else {
		return Tag{}, fmt.Errorf("Tag must have a name")
    }

	var description *string
	if v, ok := strs["description"]; ok {
		description = &v
	}

	return Tag{
		name:            *name,
		description:     description,
	}, nil
}
