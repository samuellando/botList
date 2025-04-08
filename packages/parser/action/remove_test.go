package action

import (
	"fedilist/packages/parser/jsonld"
	"testing"
	"time"
)

func TestParseRemove(t *testing.T) {
	json := `
    {
        "@context": [
            "https://schema.org",
            {
                "owner": "https://fedilist.com/owner",
                "editor": "https://fedilist.com/editor",
                "viewer": "https://fedilist.com/viewer",
                "atIndex": "https://fedilist.com/atIndex",
                "Result": "https://fedilist.com/Result"
            }
        ],
        "@type": "RemoveAction",
        "agent": {
            "@type": "Person",
            "@id": "fedilist.com/people/sam",
            "name": "Samuel Lando",
            "description": "Founder"
        },
        "object": {
            "@type": "ItemList",
            "name": "Bee Movie"
        },
        "atIndex": 25,
        "targetCollection": {
            "@id": "fedilist.com/list/2",
            "@type": "ItemList"
        },
        "startTime": "2025-04-02T10:30:00Z",
        "endTime": "2012-04-24T18:25:43Z",
        "result": {
            "@type": "Result",
            "identifier": "200",
            "description": "Success"
        }
    }
    `
	raw, err := jsonld.Expand([]byte(json))
	if err != nil {
		panic(err)
	}
	anyA, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	switch a := anyA.(type) {
	case Remove:
		if a.Agent().Id != "fedilist.com/people/sam" {
			t.Fatal("Agent did not load")
		}
		if *a.Object().Name() != "Bee Movie" {
			t.Fatal("Object did not load")
		}
		s := "2025-04-02T10:30:00Z"
		ti, _ := time.Parse(time.RFC3339, s)
		if !a.StartTime().Equal(ti) {
			t.Fatalf("startTime did not load %s", a.StartTime().Sub(ti))
		}
		s = "2012-04-24T18:25:43Z"
		ti, _ = time.Parse(time.RFC3339, s)
		if !ti.Equal(*a.EndTime()) {
			t.Fatalf("endTime did not load, %s", a.EndTime().Sub(ti))
		}
		if a.Result().Identifier != "200" {
			t.Fatal("result did not load")
		}
        if a.AtIndex() != 25 {
			t.Fatalf("index did not load %d", a.AtIndex())
        }
    default:
        t.Fatal("Wrong type")
	}
}
