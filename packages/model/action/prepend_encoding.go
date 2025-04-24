package action

import (
	"fedilist/packages/jsonld"
	"fmt"
)


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
