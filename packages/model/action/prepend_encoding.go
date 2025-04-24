package action

import (
	"encoding/json"
	"fedilist/packages/jsonld"
	"fmt"
)

func (a Prepend) MarshalJSON() ([]byte, error) {
    tl := a.targetListAction.marshal()
    tl.Type = "http://schema.org/PrependAction"
	return json.Marshal(tl)
}


func parsePrepend(json map[string]any) (Prepend, error) {
	if jsonld.GetType(json) != "http://schema.org/PrependAction" {
		return Prepend{}, fmt.Errorf("Wrong @type")
	}
	tla, err := parseTargetListAction(json)
	if err != nil {
		return Prepend{}, err
	}
	return Prepend{targetListAction: tla}, nil
}
