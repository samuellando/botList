package action

import (
	"fedilist/packages/parser/jsonld"
	"fmt"
)

var CONTEXT = map[string]any{
	"@context": []any{
		"http://schema.org",
		map[string]any{
			"owner":   "https://fedilist.org/owner",
			"editor":  "https://fedilist.org/editor",
			"viewer":  "https://fedilist.org/viewer",
			"atIndex": "https://fedilist.org/toIndex",
			"Result":  "https://fedilist.com/Result",
		},
	},
}

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
	default:
		return nil, fmt.Errorf("Unrecognized action")
	}
}
