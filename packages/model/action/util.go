package action

import (
	"fedilist/packages/jsonld"
	"fmt"
)

func Parse(json map[string]any) (Action, error) {
	switch jsonld.GetType(json) {
	case "http://schema.org/CreateAction":
		return parseCreate(json)
	case "http://schema.org/AppendAction":
		return parseAppend(json)
	case "http://schema.org/PrependAction":
		return parsePrepend(json)
	case "http://schema.org/InsertAction":
		return parseInsert(json)
	case "http://schema.org/RemoveAction":
		return parseRemove(json)
	case "http://schema.org/UpdateAction":
		return parseUpdate(json)
	case "http://schema.org/DeleteAction":
		return parseDelete(json)
	case "http://fedilist.com/ExecuteAction":
		return parseExecute(json)
	default:
		return nil, fmt.Errorf("Unrecognized action")
	}
}
