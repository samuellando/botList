package list

import (
	"fedilist/packages/model/list"
	"strconv"
)

type ListStore struct {
	base string
	db   map[string]list.ItemList
	keys map[string][]byte
}

func CreateStore(base string) ListStore {
	return ListStore{
		base: base,
		keys: make(map[string][]byte),
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
		ilv.Id = id
		ilv.Name = l.Name()
		ilv.Description = l.Description()
		ilv.Url = l.Url()
		ilv.Tags = l.Tags()
		ilv.ItemListElement = l.ItemListElement()
		ilv.Hooks = l.Hooks()
		ilv.Key = l.Key()
	})
	s.db[id] = withId
	return withId, nil
}

func (s ListStore) Append(to, e list.ItemList) (list.ItemList, error) {
	l := s.db[to.Id()]
	l.Append(e)
	s.db[l.Id()] = l
	return l, nil
}

func (s ListStore) Prepend(to, e list.ItemList) (list.ItemList, error) {
	l := s.db[to.Id()]
	l.Prepend(e)
	s.db[l.Id()] = l
	return l, nil
}


func (s ListStore) Remove(to, e list.ItemList) (list.ItemList, error) {
	l := s.db[to.Id()]
	l.Remove(e)
	s.db[l.Id()] = l
	return l, nil
}

func (s ListStore) Update(to, e list.ItemList) (list.ItemList, error) {
	withId := list.Create(func(ilv *list.ItemListValues) {
		ilv.Id = to.Id()
		ilv.Name = e.Name()
		ilv.Description = e.Description()
		ilv.Url = e.Url()
		ilv.Tags = e.Tags()
		ilv.ItemListElement = e.ItemListElement()
		ilv.Hooks = e.Hooks()
		ilv.Key = e.Key()
	})
	s.db[to.Id()] = withId
	return withId, nil
}

func (s ListStore) StoreKey(l list.ItemList, key []byte) error {
    s.keys[l.Id()] = key
    return nil
}

func (s ListStore) GetKey(l list.ItemList) ([]byte, error) {
    return s.keys[l.Id()], nil
}

