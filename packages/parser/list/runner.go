package list

import "fmt"

type Runner struct {
    id string
    name string
    inbox string
}

type RunnerValues struct {
    Id string
    Name string
    Inbox string
}

func CreateRunner(fs ...func(*RunnerValues)) (Runner, error) {
	v := RunnerValues{}
	for _, f := range fs {
		f(&v)
	}
    if v.Name == "" || v.Name == "" || v.Inbox == "" {
        return Runner{}, fmt.Errorf("Runner requies name, id and inbox URL")
    }
	return Runner{
        id: v.Id,
        name: v.Name,
        inbox: v.Inbox,
	}, nil
}
