package action

import (
	"fedilist/packages/jsonld"
	"fmt"
    "encoding/json"
)

func (a Remove) MarshalJSON() ([]byte, error) {
    tl := a.targetListAction.marshal()
    tl.Type = "http://schema.org/RemoveAction"
	return json.Marshal(tl)
}

func parseRemove(json map[string]any) (Remove, error) {
	if jsonld.GetType(json) != "http://schema.org/RemoveAction" {
		return Remove{}, fmt.Errorf("Wrong @type")
	}
	tla, err := parseTargetListAction(json)
	if err != nil {
		return Remove{}, err
	}
	return Remove{targetListAction: tla}, nil
}
