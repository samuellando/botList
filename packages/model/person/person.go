package person

import (
	"fedilist/packages/model/list"
)

type Person struct {
	id          string
	name        string
	description string
	list        []list.ItemList
	key         string
	inbox       string
	outbox      string
}

func (p Person) Id() string {
	return p.id
}

func (p Person) Name() string {
	return p.name
}

func (p Person) Description() string {
	return p.description
}

func (p Person) List() []list.ItemList {
	return p.list
}

func (p Person) Key() string {
	return p.key
}

func (p Person) Inbox() string {
	return p.inbox
}

func (p Person) Outbox() string {
	return p.outbox
}

type PersonValues struct {
	Id          string
	Name        string
	Description string
	Key         string
	List        []list.ItemList
	Inbox       string
	Outbox      string
}

func CreatePerson(fs ...func(*PersonValues)) Person {
	pv := PersonValues{
		List: make([]list.ItemList, 0),
	}
	for _, f := range fs {
		f(&pv)
	}
	p := Person{
		id:          pv.Id,
		name:        pv.Name,
		description: pv.Description,
		list:        pv.List,
		key:         pv.Key,
		inbox:       pv.Inbox,
		outbox:      pv.Outbox,
	}
	return p
}

func (p *Person) AddList(l list.ItemList) {
	p.list = append(p.list, l)
}
