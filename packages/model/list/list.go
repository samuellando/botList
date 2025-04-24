package list

import (
	"fedilist/packages/model/hook"
	"fedilist/packages/model/tag"
	"slices"
)

type ItemList struct {
	id              string
	name            string
	description     string
	url             string
	tags            []tag.Tag
	hooks           []hook.Hook
	numberOfItems   int
	itemListElement []ItemList
	key             string
}

func (l ItemList) Id() string {
	return l.id
}

func (l ItemList) Name() string {
	return l.name
}

func (l ItemList) Description() string {
	return l.description
}

func (l ItemList) Url() string {
	return l.url
}

func (l ItemList) Tags() []tag.Tag {
	return l.tags
}

func (l ItemList) Hooks() []hook.Hook {
	return l.hooks
}

func (l ItemList) Key() string {
	return l.key
}

func (l ItemList) ItemListElement() []ItemList {
	return l.itemListElement
}

type ItemListValues struct {
	Id              string
	Name            string
	Description     string
	Url             string
	Tags            []tag.Tag
	Hooks           []hook.Hook
	ItemListElement []ItemList
	Key             string
}

func Create(fs ...func(*ItemListValues)) ItemList {
	p := ItemListValues{
		Tags:            make([]tag.Tag, 0),
		Hooks:           make([]hook.Hook, 0),
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
		numberOfItems:   n,
		itemListElement: p.ItemListElement,
		tags:            p.Tags,
		hooks:           p.Hooks,
		key:             p.Key,
	}
}

func (l *ItemList) Append(e ItemList) {
	l.itemListElement = append(l.itemListElement, e)
	n := len(l.itemListElement)
	l.numberOfItems = n
}	

func (l *ItemList) Prepend(e ItemList) {
	l.itemListElement = slices.Insert(l.itemListElement, 0, e)
	n := len(l.itemListElement)
	l.numberOfItems = n
}	

func (l *ItemList) Remove(e ItemList) {
	l.itemListElement = slices.DeleteFunc(l.itemListElement, func(le ItemList) bool {
		return le.Id() == e.Id()
	})
	n := len(l.itemListElement)
	l.numberOfItems = n
}
