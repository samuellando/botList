package list

import (
	"fedilist/packages/model/list"
	"strconv"
	"encoding/json"
	"os"
)

type ListStore struct {
	base string
	db   map[string]list.ItemList
	keys map[string][]byte
}

func (ps ListStore) MarshalJSON() ([]byte, error) {
	type Alias ListStore
	return json.Marshal(&struct {
		Base string                   `json:"base"`
		Db   map[string]list.ItemList `json:"db"`
		Keys map[string][]byte        `json:"keys"`
		*Alias
	}{
		Base:  ps.base,
		Db:    ps.db,
		Keys:  ps.keys,
		Alias: (*Alias)(&ps),
	})
}

func (ps *ListStore) UnmarshalJSON(data []byte) error {
	type Alias ListStore
	aux := &struct {
		Base string                   `json:"base"`
		Db   map[string]list.ItemList `json:"db"`
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

func save(ps ListStore) {
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
	err = os.WriteFile(dirPath+"/lists.json", b, 0644)
	if err != nil {
		panic(err)
	}
}

func load() (ListStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ListStore{}, err
	}
	filePath := homeDir + "/.cache/botlist/lists.json"
	b, err := os.ReadFile(filePath)
	if err != nil {
		return ListStore{}, err
	}
	var ps ListStore
	err = json.Unmarshal(b, &ps)
	if err != nil {
		return ListStore{}, err
	}
	return ps, nil
}

func CreateStore(base string) ListStore {
	ps, err := load()
	if err == nil {
		return ps
	}
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
	save(s)
	return withId, nil
}

func (s ListStore) Append(to, e list.ItemList) (list.ItemList, error) {
	l := s.db[to.Id()]
	l.Append(e)
	s.db[l.Id()] = l
	save(s)
	return l, nil
}

func (s ListStore) Prepend(to, e list.ItemList) (list.ItemList, error) {
	l := s.db[to.Id()]
	l.Prepend(e)
	s.db[l.Id()] = l
	save(s)
	return l, nil
}


func (s ListStore) Remove(to, e list.ItemList) (list.ItemList, error) {
	l := s.db[to.Id()]
	l.Remove(e)
	s.db[l.Id()] = l
	save(s)
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
	save(s)
	return withId, nil
}

func (s ListStore) StoreKey(l list.ItemList, key []byte) error {
    s.keys[l.Id()] = key
	save(s)
    return nil
}

func (s ListStore) GetKey(l list.ItemList) ([]byte, error) {
    return s.keys[l.Id()], nil
}

