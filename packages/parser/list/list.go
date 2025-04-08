package list

import (
	"fedilist/packages/parser/jsonld"
)

type ItemList struct {
	id              *string
	name            *string
	description     *string
	url             *string
	tags            []Tag
	numberOfItems   *int
	itemListElement []ItemList
}

func (l ItemList) Id() *string {
	return l.id
}

func (l ItemList) Name() *string {
	return l.name
}

func (l ItemList) Description() *string {
	return l.description
}

func (l ItemList) Url() *string {
	return l.url
}

func (l ItemList) Tags() []Tag {
	return l.tags
}

func (l ItemList) ItemListElement() []ItemList {
	return l.itemListElement
}

type ItemListParam func(*ItemListValues)

type ItemListValues struct {
	Id              *string
	Name            *string
	Description     *string
	Url             *string
	Tags            []Tag
	ItemListElement []ItemList
}

func Create(fs ...ItemListParam) ItemList {
	p := ItemListValues{
		Tags:            make([]Tag, 0),
		ItemListElement: make([]ItemList, 0),
	}
	for _, f := range fs {
		f(&p)
	}
	n := len(p.ItemListElement)
	return ItemList{
		id:              p.Id,
		name:            p.Name,
		description:     p.Description,
		numberOfItems:   &n,
		itemListElement: p.ItemListElement,
		tags:            p.Tags,
	}
}

func (l *ItemList) Append(e ItemList) {
	l.itemListElement = append(l.itemListElement, e)
	n := len(l.itemListElement)
	l.numberOfItems = &n
}

func (l ItemList) Serialize() []byte {
	return jsonld.Marshal(CONTEXT, l)
}
