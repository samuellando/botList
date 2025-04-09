package action

import (
	"fedilist/packages/parser/jsonld"
	"fmt"
    "encoding/json"
)

func (a Append) MarshalJSON() ([]byte, error) {
    tl := a.targetListAction.marshal()
    tl.Type = "http://schema.org/AppendAction"
	return json.Marshal(tl)
}

func parseAppend(json map[string]any) (Append, error) {
	if jsonld.GetType(json) != "http://schema.org/AppendAction" {
		return Append{}, fmt.Errorf("Wrong @type")
	}
	tla, err := parseTargetListAction(json)
	if err != nil {
		return Append{}, err
	}
	return Append{targetListAction: tla}, nil
}
