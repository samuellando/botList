package util

import (
	"fedilist/packages/parser/jsonld"
	"io"
	"net/http"
)

func GetBodyJsonld(req *http.Request) (map[string]any, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	data, err := jsonld.Expand(bodyBytes)
	return data, err
}
