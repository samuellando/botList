package action

import (
	"fedilist/packages/jsonld"
	"fmt"
)

func parseInsert(json map[string]any) (Insert, error) {
	if jsonld.GetType(json) != "http://schema.org/InsertAction" {
		return Insert{}, fmt.Errorf("Wrong @type")
	}
	tla, err := parseTargetListAction(json)
	if err != nil {
		return Insert{}, err
	}
	schemaOrgValues := jsonld.GetNamespaceValues(json, "http://fedilist.com")
	ints := jsonld.GetBaseTypeValues[float64](schemaOrgValues)
	var atIndex int
	if i, ok := ints["atIndex"]; ok {
		atIndex = int(i)
	}
	return Insert{targetListAction: tla, atIndex: atIndex}, nil
}
