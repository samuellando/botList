package util

import (
	"fedilist/packages/jsonld"
	"encoding/json"
	"io"
	"crypto/ed25519"
	"encoding/base64"
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

type Signable[T any] interface {
	Sign(string) T
}

func GetSignature[T any](o Signable[T], seed []byte) (string, error) {
	privateKey := ed25519.NewKeyFromSeed(seed)
	noSig := o.Sign("")
	txt, err := json.Marshal(noSig)
	if err != nil {
		return "", err
	}
	sig := ed25519.Sign(privateKey, txt)
	return base64.StdEncoding.EncodeToString(sig), nil
}
