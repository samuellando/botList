package action

import (
	"fedilist/packages/jsonld"
	"fmt"
)

func parseUpdate(json map[string]any) (Update, error) {
	if jsonld.GetType(json) != "http://schema.org/UpdateAction" {
		return Update{}, fmt.Errorf("Wrong @type")
	}
	tla, err := parseTargetListAction(json)
	if err != nil {
		return Update{}, err
	}
	return Update{targetListAction: tla}, nil
}
