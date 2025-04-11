package list

import "fmt"

type Runner struct {
	id      string
	name    string
	inbox   string
	service []Service
}

type Service struct {
	Name   string
	Schema string
}

type RunnerValues struct {
	Id      string
	Name    string
	Inbox   string
	Service []Service
}

func (r Runner) Id() *string {
	return &r.id
}

func (r Runner) Name() *string {
	return &r.name
}

func (r Runner) Inbox() *string {
	return &r.inbox
}

func (r Runner) Services() []Service {
	return r.service
}

func CreateRunner(fs ...func(*RunnerValues)) (Runner, error) {
	v := RunnerValues{}
	for _, f := range fs {
		f(&v)
	}
	if v.Id == "" || v.Name == "" || v.Inbox == "" {
		return Runner{}, fmt.Errorf("Runner requies name, id and inbox URL")
	}
	return Runner{
		id:    v.Id,
		name:  v.Name,
		inbox: v.Inbox,
        service: v.Service,
	}, nil
}
