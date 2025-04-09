package list

import "fmt"

type Tag struct {
	name  string
	description *string
}

type TagValues struct {
	Name  string
	Description *string
}

func CreateTag(fs ...func(*TagValues)) (Tag, error) {
    v := TagValues{}
    for _, f := range fs {
        f(&v)
    }
    if v.Name == "" {
        return Tag{}, fmt.Errorf("Tag name must be provided")
    }
    return Tag{
        name: v.Name,
        description: v.Description,
    }, nil
}

func (t Tag) Name() string {
    return t.name
}

func (t Tag) Description() *string {
    return t.description
}
