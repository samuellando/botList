package list

import (
	"encoding/json"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/hook"
	"fedilist/packages/model/tag"
	"fmt"
)

func (l ItemList) MarshalJSON() ([]byte, error) {
	type External struct {
		Type            string      `json:"@type"`
		Id              string     `json:"@id,omitempty"`
		Name            string     `json:"http://schema.org/name,omitempty"`
		Description     string     `json:"http://schema.org/description,omitempty"`
		Url             string     `json:"http://schema.org/url,omitempty"`
		Tags            []tag.Tag   `json:"http://schema.org/tags,omitempty"`
		Hooks           []hook.Hook `json:"http://fedilist.com/hooks,omitempty"`
		NumberOfItems   int        `json:"http://schema.org/numberOfItems"`
		ItemListElement []ItemList  `json:"http://schema.org/itemListElement,omitempty"`
		Key             string      `json:"http://fedilist.com/key,omitempty"`
	}
	return json.Marshal(External{
		Type:            "http://schema.org/ItemList",
		Id:              l.id,
		Name:            l.name,
		Description:     l.description,
		Url:             l.url,
		Tags:            l.tags,
		Hooks:           l.hooks,
		NumberOfItems:   l.numberOfItems,
		ItemListElement: l.itemListElement,
		Key:             l.key,
	})
}

func (l *ItemList) UnmarshalJSON(data []byte) error {
	type External struct {
		Type            string      `json:"@type"`
		Id              string      `json:"@id,omitempty"`
		Name            string      `json:"http://schema.org/name,omitempty"`
		Description     string      `json:"http://schema.org/description,omitempty"`
		Url             string      `json:"http://schema.org/url,omitempty"`
		Tags            []tag.Tag   `json:"http://schema.org/tags,omitempty"`
		Hooks           []hook.Hook `json:"http://fedilist.com/hooks,omitempty"`
		NumberOfItems   int         `json:"http://schema.org/numberOfItems"`
		ItemListElement []ItemList  `json:"http://schema.org/itemListElement,omitempty"`
		Key             string      `json:"http://fedilist.com/key,omitempty"`
	}
	var ext External
	if err := json.Unmarshal(data, &ext); err != nil {
		return err
	}
	if ext.Type != "http://schema.org/ItemList" {
		return fmt.Errorf("Type must be ItemList")
	}
	l.id = ext.Id
	l.name = ext.Name
	l.description = ext.Description
	l.url = ext.Url
	l.tags = ext.Tags
	l.hooks = ext.Hooks
	l.numberOfItems = ext.NumberOfItems
	l.itemListElement = ext.ItemListElement
	l.key = ext.Key
	return nil
}

func Parse(json map[string]any) (ItemList, error) {
	if jsonld.GetType(json) != "http://schema.org/ItemList" {
		return ItemList{}, fmt.Errorf("Type must be ItemList")
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	fediOrgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)
	ints := jsonld.GetBaseTypeValues[float64](schemaOrgValues)
	objLs := jsonld.GetCompositeTypeArrayValues(schemaOrgValues)
	var name string
	if v, ok := strs["name"]; ok {
		name = v
	}
	var numberOfItems int
	if v, ok := ints["numberOfItems"]; ok {
		numberOfItems = int(v)
	}

	var id string = jsonld.GetId(json)
	var description string
	if v, ok := strs["description"]; ok {
		description = v
	}
	var url string
	if v, ok := strs["url"]; ok {
		url = v
	}

	var tags []tag.Tag
	if l, ok := objLs["tags"]; ok {
		tags = make([]tag.Tag, len(l))
		for i, v := range l {
			tag, err := tag.Parse(v)
			if err != nil {
				return ItemList{}, err
			}
			tags[i] = tag
		}
	}

	var itemListElement []ItemList
	if l, ok := objLs["itemListElement"]; ok {
		itemListElement = make([]ItemList, len(l))
		for i, v := range l {
			elem, err := Parse(v)
			if err != nil {
				return ItemList{}, err
			}
			itemListElement[i] = elem
		}
	}

	objLs = jsonld.GetCompositeTypeArrayValues(fediOrgValues)
	var hooks []hook.Hook
	if l, ok := objLs["hooks"]; ok {
		hooks = make([]hook.Hook, len(l))
		for i, v := range l {
			elem, err := hook.Parse(v)
			if err != nil {
				return ItemList{}, err
			}
			hooks[i] = elem
		}
	}

	strs = jsonld.GetBaseTypeValues[string](fediOrgValues)

	var key string
	if v, ok := strs["key"]; ok {
		key = v
	}

	return ItemList{
		id:              id,
		name:            name,
		description:     description,
		url:             url,
		tags:            tags,
		hooks:           hooks,
		numberOfItems:   numberOfItems,
		itemListElement: itemListElement,
		key:             key,
	}, nil
}
