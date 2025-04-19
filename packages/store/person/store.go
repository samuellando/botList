package person

import (
	"fedilist/packages/model/list"
	"fedilist/packages/model/person"
	"fmt"
	"strings"
)

type PersonStore struct {
	base string
	db   map[string]person.Person
    keys map[string][]byte
}

func CreateStore(base string) PersonStore {
	return PersonStore{
		base: base,
		db:   make(map[string]person.Person),
		keys:   make(map[string][]byte),
	}
}

func (s PersonStore) GetById(id string) (person.Person, error) {
	if p, ok := s.db[strings.ToLower(id)]; ok {
		return p, nil
	} else {
		return person.Person{}, fmt.Errorf("not found")
	}
}

func (s PersonStore) StoreKey(p person.Person, key []byte) error {
    s.keys[p.Id()] = key
    return nil
}

func (s PersonStore) GetKey(p person.Person) ([]byte, error) {
    return s.keys[p.Id()], nil
}

func (s PersonStore) GetByPartialId(pid string) (person.Person, error) {
	return s.GetById(s.base + pid)
}

func (s PersonStore) Insert(p person.Person) (person.Person, error) {
	if p.Name() == "" {
		panic("Empty name")
	}
	id := strings.ToLower(s.base + p.Name())
	if _, ok := s.db[id]; ok {
		return person.Person{}, fmt.Errorf("Already exists")
	}
	p = person.CreatePerson(func(pv *person.PersonValues) {
		pv.Id = id
		pv.Name = p.Name()
		pv.Description = p.Description()
		pv.List = p.List()
        pv.Key = p.Key()
        pv.Inbox = strings.ToLower(s.base + p.Name() + "/inbox")
        pv.Outbox = strings.ToLower(s.base + p.Name() + "/outbox")
	})
	s.db[id] = p
	return p, nil
}

func (s PersonStore) AddList(p person.Person, l list.ItemList) (person.Person, error) {
	dbp, ok := s.db[p.Id()]
	if !ok {
		return person.Person{}, fmt.Errorf("Not found")
	}
	dbp.AddList(l)
	s.db[p.Id()] = dbp
	return dbp, nil
}
