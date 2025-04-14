package list

import (
	"fedilist/packages/model/list"
	"strconv"
)

type ListStore struct {
	base string
	db   map[string]list.ItemList
}

func CreateStore(base string) ListStore {
	return ListStore{
		base: base,
		db:   make(map[string]list.ItemList),
	}
}

func (s ListStore) GetById(id string) (list.ItemList, error) {
	return s.db[id], nil
}

func (s ListStore) GetByPartialId(id string) (list.ItemList, error) {
	return s.db[s.base+id], nil
}

func (s ListStore) Insert(l list.ItemList) (list.ItemList, error) {
	id := string(s.base + strconv.Itoa(len(s.db)))
	withId := list.Create(func(ilv *list.ItemListValues) {
		ilv.Id = &id
		ilv.Name = l.Name()
		ilv.Description = l.Description()
		ilv.Url = l.Url()
		ilv.Tags = l.Tags()
		ilv.ItemListElement = l.ItemListElement()
		ilv.Hooks = l.Hooks()
	})
	s.db[id] = withId
	return withId, nil
}

func (s ListStore) Append(to, e list.ItemList) (list.ItemList, error) {
	l := s.db[*to.Id()]
	l.Append(e)
	s.db[*l.Id()] = l
	return l, nil
}
