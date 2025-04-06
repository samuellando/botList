package list

import (
	"encoding/json"
	"fedilist/packages/parser/jsonld"
	"fmt"
	"strconv"
)

type ItemList struct {
	Id              *string
	Name            *string
	Description     *string
	Url             *string
	Tags            []Tag
	NumberOfItems   *int
	ItemListElement []ItemList
}

func (l ItemList) MarshalJSON() ([]byte, error) {
	type External struct {
        Type            string     `json:"@type"`
		Id              *string    `json:"@id,omitempty"`
		Name            *string    `json:"http://schema.org/name,omitempty"`
		Description     *string    `json:"http://schema.org/description,omitempty"`
		Url             *string    `json:"http://schema.org/url,omitempty"`
		Tags            []Tag      `json:"http://schema.org/tags,omitempty"`
		NumberOfItems   *int       `json:"http://schema.org/numberOfItems,omitempty"`
		ItemListElement []ItemList `json:"http://schema.org/itemListElement,omitempty"`
	}
	return json.Marshal(External{
		Type:            "http://schema.org/ItemList",
		Id:              l.Id,
		Name:            l.Name,
		Description:     l.Description,
		Url:             l.Url,
		Tags:            l.Tags,
		NumberOfItems:   l.NumberOfItems,
		ItemListElement: l.ItemListElement,
	})
}

type Tag struct {
	Identifier  string
	Description string
}

var DB = make(map[string]ItemList)

func CreateList(domain, name, description string) ItemList {
	id := string(domain + "/list/" + strconv.Itoa(len(DB)))
    n := 0
	l := ItemList{
		Id:          &id,
		Name:        &name,
		Description: &description,
        NumberOfItems: &n,
        ItemListElement: make([]ItemList, 0),
	}
	DB[*l.Id] = l
	return l
}

func GetListById(id string) ItemList {
    return DB[id]
}

func (l *ItemList) Append(e ItemList) {
    l.ItemListElement = append(l.ItemListElement, e)
    n := *l.NumberOfItems + 1
    l.NumberOfItems = &n
    DB[*l.Id] = *l
}

func LoadItemList(json map[string]any) (ItemList, error) {
	if jsonld.GetType(json) != "http://schema.org/ItemList" {
		return ItemList{}, fmt.Errorf("Type must be ItemList")
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
	strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)
	ints := jsonld.GetBaseTypeValues[int](schemaOrgValues)
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
	var itemListElement []ItemList

	return ItemList{
		Id:              id,
		Name:            name,
		Description:     description,
		Url:             url,
		Tags:            tags,
		NumberOfItems:   numberOfItems,
		ItemListElement: itemListElement,
	}, nil
}
