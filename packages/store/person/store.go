package person

import (
	"encoding/json"
	"fedilist/packages/model/list"
	"fedilist/packages/model/person"
	"fmt"
	"os"
	"strings"
)

type PersonStore struct {
	base string
	db   map[string]person.Person
	keys map[string][]byte
}

func (ps PersonStore) MarshalJSON() ([]byte, error) {
	type Alias PersonStore
	return json.Marshal(&struct {
		Base string                   `json:"base"`
		Db   map[string]person.Person `json:"db"`
		Keys map[string][]byte        `json:"keys"`
		*Alias
	}{
		Base:  ps.base,
		Db:    ps.db,
		Keys:  ps.keys,
		Alias: (*Alias)(&ps),
	})
}

func (ps *PersonStore) UnmarshalJSON(data []byte) error {
	type Alias PersonStore
	aux := &struct {
		Base string                   `json:"base"`
		Db   map[string]person.Person `json:"db"`
		Keys map[string][]byte        `json:"keys"`
		*Alias
	}{
		Alias: (*Alias)(ps),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ps.base = aux.Base
	ps.db = aux.Db
	ps.keys = aux.Keys
	return nil
}

func save(ps PersonStore) {
	b, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		panic(err)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dirPath := homeDir + "/.cache/botlist"
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(dirPath+"/users.json", b, 0644)
	if err != nil {
		panic(err)
	}
}

func load() (PersonStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return PersonStore{}, err
	}
	filePath := homeDir + "/.cache/botlist/users.json"
	b, err := os.ReadFile(filePath)
	if err != nil {
		return PersonStore{}, err
	}
	var ps PersonStore
	err = json.Unmarshal(b, &ps)
	if err != nil {
		return PersonStore{}, err
	}
	return ps, nil
}

func CreateStore(base string) PersonStore {
	ps, err := load()
	if err == nil {
		return ps
	}
	return PersonStore{
		base: base,
		db:   make(map[string]person.Person),
		keys: make(map[string][]byte),
	}
}

func (s PersonStore) GetById(id string) (person.Person, error) {
	fmt.Print(">>>", id, s.db)
	if p, ok := s.db[strings.ToLower(id)]; ok {
		return p, nil
	} else {
		return person.Person{}, fmt.Errorf("not found")
	}
}

func (s PersonStore) StoreKey(p person.Person, key []byte) error {
	s.keys[p.Id()] = key
	save(s)
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
	save(s)
	return p, nil
}

func (s PersonStore) AddList(p person.Person, l list.ItemList) (person.Person, error) {
	dbp, ok := s.db[p.Id()]
	if !ok {
		return person.Person{}, fmt.Errorf("Not found")
	}
	dbp.AddList(l)
	s.db[p.Id()] = dbp
	save(s)
	return dbp, nil
}
