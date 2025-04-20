package runner

import "fmt"

type Runner struct {
	id      string
	name    string
	inbox   string
	service []Service
	key     string
}

type RunnerValues struct {
	Id      string
	Name    string
	Inbox   string
	Service []Service
	Key     string
}

func (r Runner) Id() string {
	return r.id
}

func (r Runner) Name() string {
	return r.name
}

func (r Runner) Inbox() string {
	return r.inbox
}

func (r Runner) Key() string {
	return r.key
}

func (r Runner) Services() []Service {
	return r.service
}

func Create(fs ...func(*RunnerValues)) (Runner, error) {
	v := RunnerValues{}
	for _, f := range fs {
		f(&v)
	}
	if v.Id == "" || v.Name == "" || v.Inbox == "" {
		return Runner{}, fmt.Errorf("Runner requies name, id and inbox URL")
	}
	return Runner{
		id:      v.Id,
		name:    v.Name,
		inbox:   v.Inbox,
		service: v.Service,
		key:     v.Key,
	}, nil
}
