package action

import (
	"fedilist/packages/parser/jsonld"
	"testing"
	"time"
)

func TestParseCreate(t *testing.T) {
	json := `
    {
        "@context": [
            "https://schema.org",
            {
                "owner": "http://fedilist.com/owner",
                "editor": "http://fedilist.com/editor",
                "viewer": "http://fedilist.com/viewer",
                "atIndex": "http://fedilist.com/toIndex",
                "Result": "http://fedilist.com/Result"
            }
        ],
        "@type": "CreateAction",
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
	case Create:
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
    default:
        t.Fatal("Wrong type")
	}
}
