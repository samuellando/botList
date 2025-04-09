package list

import (
	"encoding/json"
	"fedilist/packages/parser/jsonld"
	"fmt"
)

var CONTEXT = map[string]any{
	"@context": []any{
		"http://schema.org",
		map[string]any{
			"owner":   "https://fedilist.org/owner",
			"editor":  "https://fedilist.org/editor",
			"viewer":  "https://fedilist.org/viewer",
			"hooks":   "https://fedilist.org/hooks",
			"inbox":   "https://fedilist.org/inbox",
			"atIndex": "https://fedilist.org/toIndex",
			"Result":  "https://fedilist.com/Result",
		},
	},
}

func (l ItemList) MarshalJSON() ([]byte, error) {
	type External struct {
		Type            string     `json:"@type"`
		Id              *string    `json:"@id,omitempty"`
		Name            *string    `json:"http://schema.org/name,omitempty"`
		Description     *string    `json:"http://schema.org/description,omitempty"`
		Url             *string    `json:"http://schema.org/url,omitempty"`
		Tags            []Tag      `json:"http://schema.org/tags,omitempty"`
		Hooks           []Hook     `json:"http://fedilist.com/hooks,omitempty"`
		NumberOfItems   *int       `json:"http://schema.org/numberOfItems,omitempty"`
		ItemListElement []ItemList `json:"http://schema.org/itemListElement,omitempty"`
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
	})
}

func Parse(json map[string]any) (ItemList, error) {
	if jsonld.GetType(json) != "http://schema.org/ItemList" {
		return ItemList{}, fmt.Errorf("Type must be ItemList")
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)
	ints := jsonld.GetBaseTypeValues[int](schemaOrgValues)
	objLs := jsonld.GetCompositeTypeArrayValues(schemaOrgValues)
	var name *string
	if v, ok := strs["name"]; ok {
		name = &v
	}
	var numberOfItems *int
	if v, ok := ints["numberOfItems"]; ok {
		numberOfItems = &v
	}

	var id *string = jsonld.GetId(json)
	var description *string
	if v, ok := strs["description"]; ok {
		description = &v
	}
	var url *string
	if v, ok := strs["url"]; ok {
		url = &v
	}

	var tags []Tag
	if l, ok := objLs["tags"]; ok {
		tags = make([]Tag, len(l))
		for i, v := range l {
			tag, err := ParseTag(v)
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

	var hooks []Hook
	if l, ok := objLs["hooks"]; ok {
		hooks = make([]Hook, len(l))
		for i, v := range l {
			elem, err := ParseHook(v)
			if err != nil {
				return ItemList{}, err
			}
			hooks[i] = elem
		}
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
	}, nil
}
