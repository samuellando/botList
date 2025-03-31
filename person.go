package main

type Person struct {
	Id      string
	Name    string
	Summary string
	Lists   List
    Inbox   []Activity[List]
}

var PEOPLE = make(map[string]*Person)

func CreatePerson(name, summary string) Person {
    id := getId("user", name)
    l := CreateList(name+"'s lists", id)
	p := Person{
		Id:      getId("user", name),
		Name:    name,
		Summary: summary,
		Lists:   l,
        Inbox: make([]Activity[List], 0),
	}
	PEOPLE[p.Id] = &p
	return p
}

func GetPersonById(id string) Person {
	if person, ok := PEOPLE[id]; ok {
		return *person
	} else {
        p := Person{
			Id:      id,
			Name:    "",
			Summary: "",
			Lists:   CreateList("", id),
            Inbox: make([]Activity[List], 0),
		}
        PEOPLE[p.Id] = &p
        return p
	}
}

func (p *Person) AddToInbox(a Activity[List]) {
    p.Inbox = append(p.Inbox, a)
    PEOPLE[p.Id].Inbox = append(PEOPLE[p.Id].Inbox, a)
}
