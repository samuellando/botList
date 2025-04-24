package action

import (
	"fedilist/packages/jsonld"
	"fmt"
)


func parseDelete(json map[string]any) (Delete, error) {
    if jsonld.GetType(json) != "http://schema.org/DeleteAction" {
        return Delete{}, fmt.Errorf("Wrong @type")
    }
    tla, err := parseTargetListAction(json)
	if err != nil {
		return Delete{}, err
	}
    return Delete{targetListAction: tla}, nil
}
