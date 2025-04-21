package jsonld

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/piprate/json-gold/ld"
)

var CONTEXT = map[string]any{
	"@context": []any{
		"http://schema.org",
		map[string]any{
			"owner":   "http://fedilist.com/owner",
			"editor":  "http://fedilist.com/editor",
			"viewer":  "http://fedilist.com/viewer",
			"atIndex": "http://fedilist.com/toIndex",
			"Result":  "http://fedilist.com/Result",
			"hooks":   "http://fedilist.com/hooks",
			"inbox":   "http://fedilist.com/inbox",
		},
	},
}

func Marshal(data any) []byte {
	s, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	var raw map[string]any
	err = json.Unmarshal(s, &raw)
	if err != nil {
		panic(err)
	}
	p := ld.NewJsonLdProcessor()
	compact, err := p.Compact(raw, CONTEXT, nil)
	if err != nil {
		panic(err)
	}
	s, err = json.Marshal(compact)
	if err != nil {
		panic(err)
	}
	return s
}

func MarshalIndent(data any) []byte {
	s, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	var raw map[string]any
	err = json.Unmarshal(s, &raw)
	if err != nil {
		panic(err)
	}
	p := ld.NewJsonLdProcessor()
	compact, err := p.Compact(raw, CONTEXT, nil)
	if err != nil {
		panic(err)
	}
	s, err = json.MarshalIndent(compact, "", "    ")
	if err != nil {
		panic(err)
	}
	return s
}

func Expand(data []byte) (map[string]any, error) {
	proc := ld.NewJsonLdProcessor()
	compact := make(map[string]any)
	err := json.Unmarshal(data, &compact)
	if err != nil {
		return nil, err
	}
	expanded, err := proc.Expand(compact, ld.NewJsonLdOptions(""))
	if err != nil {
		return nil, err
	}
	if len(expanded) == 0 {
		return nil, fmt.Errorf("Zero objects parsed")
	}
	if m, ok := expanded[0].(map[string]any); ok {
		return m, nil
	} else {
		return nil, fmt.Errorf("Could not convert to map of strings")
	}
}

func GetBaseTypeValues[T any](json map[string]any) map[string]T {
	baseTypeValues := make(map[string]T)
	for k, v := range json {
		if a, ok := v.([]any); ok {
			if len(a) == 0 {
				continue
			}
			if o, ok := a[0].(map[string]any); ok {
				if ov, ok := o["@value"]; ok {
					if t, ok := ov.(T); ok {
						baseTypeValues[k] = t
					}
				}
			}
		}
	}
	return baseTypeValues
}

func GetBaseTypeArrayValues[T any](json map[string]any) map[string][]T {
	baseTypeValues := make(map[string][]T)
	for k, v := range json {
		if a, ok := v.([]any); ok {
			for _, v := range a {
				if o, ok := v.(map[string]any); ok {
					if ov, ok := o["@value"]; ok {
						if t, ok := ov.(T); ok {
							if _, ok := baseTypeValues[k]; ok {
								baseTypeValues[k] = append(baseTypeValues[k], t)
							} else {
								baseTypeValues[k] = make([]T, 1)
								baseTypeValues[k][0] = t
							}
						}
					}
				}

			}
		}
	}
	return baseTypeValues
}

func GetCompositeTypeValues(json map[string]any) map[string]map[string]any {
	values := make(map[string]map[string]any)
	for k, v := range json {
		if a, ok := v.([]any); ok {
			if len(a) == 0 {
				continue
			}
			if o, ok := a[0].(map[string]any); ok {
				if _, ok := o["@type"]; ok {
					values[k] = o
				}
			}
		}
	}
	return values
}

func GetCompositeTypeArrayValues(json map[string]any) map[string][]map[string]any {
	values := make(map[string][]map[string]any)
	for k, v := range json {
		if a, ok := v.([]any); ok {
			if len(a) == 0 {
				continue
			}
			if o, ok := a[0].(map[string]any); ok {
				if _, ok := o["@type"]; ok {
					values[k] = make([]map[string]any, len(a))
					for i, v := range a {
						values[k][i] = v.(map[string]any)
					}
				}
			}
		}
	}
	return values
}

func GetNamespaceValues(json map[string]any, namespace string) map[string]any {
	prefix := namespace + "/"
	nsValues := make(map[string]any)
	for k, v := range json {
		if strings.HasPrefix(k, prefix) {
			nsKey := k[len(prefix):]
			nsValues[nsKey] = v
		}
	}
	return nsValues
}

func GetId(json map[string]any) string {
	if v, ok := json["@id"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func GetType(json map[string]any) string {
	if v, ok := json["@type"]; ok {
		if a, ok := v.([]any); ok {
			if s, ok := a[0].(string); ok {
				return s
			}
		}
	}
	return ""
}
