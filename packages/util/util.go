package util

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fedilist/packages/jsonld"
	"io"
)

func GetBodyJsonld(rc io.ReadCloser) (map[string]any, error) {
	bodyBytes, err := io.ReadAll(rc)
	defer rc.Close()
	if err != nil {
		return nil, err
	}
	data, err := jsonld.Expand(bodyBytes)
	return data, err
}

type Signable[T any] interface {
	Sign(string) T
	Signature() string
}

func Sign[T any](o Signable[T], seed []byte) T {
	privateKey := ed25519.NewKeyFromSeed(seed)
	// Clear any existing signature
	noSig := o.Sign("")
	// Get a signature for the message
	txt, err := json.Marshal(noSig)
	if err != nil {
		panic(err)
	}
	signature := base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, txt))
	return o.Sign(signature)
}

func VerifySignature[T any](o Signable[T], key string) (bool, error) {
	noSig := o.Sign("")
	txt, err := json.Marshal(noSig)
	if err != nil {
		return false, err
	}
	// Decode the base64-encoded public key
	publicKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return false, err
	}
	// Retrieve the signature from the object
	signature, err := base64.StdEncoding.DecodeString(o.Signature())
	if err != nil {
		return false, err
	}
	// Verify the signature
	isValid := ed25519.Verify(publicKey, txt, signature)
	return isValid, nil
}
