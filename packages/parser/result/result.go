package result

import (
	"fmt"
    "fedilist/packages/parser/jsonld"
    "encoding/json"
)

type Result struct {
	Identifier  string
	Description string
}

func (a Result) MarshalJSON() ([]byte, error) {
	type External struct {
        Type string `json:"@type"`
        Identifier string `json:"http://schema.org/identifier"`
        Description string `json:"http://schema.org/description"`
	}
	return json.Marshal(External{
        Type: "http://fedilist.com/Result",
        Identifier: a.Identifier,
        Description: a.Description,
	})
}

func Create(identifier, description string) Result {
    return Result{
        Identifier: identifier,
        Description: description,
    }
}


func LoadResult(json map[string]any) (Result, error) {
    if jsonld.GetType(json) != "http://fedilist.com/Result" {
        return Result{}, fmt.Errorf("Cannot load non result")
    }
    schemaOrgValues := jsonld.GetNamespaceValues(json, "http://schema.org")
    strs := jsonld.GetBaseTypeValues[string](schemaOrgValues)
    
    var identifier string
    if v, ok := strs["identifier"]; ok {
        identifier = v
    } else {
        return Result{}, fmt.Errorf("Result must have a identifier")
    }
    var description string
    if v, ok := strs["description"]; ok {
        description = v
    } else {
        return Result{}, fmt.Errorf("Result must have a description")
    }

    return Result{
        Identifier: identifier,
        Description: description,
    }, nil
}
