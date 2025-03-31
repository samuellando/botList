package main

import "strconv"

// Links?
// General response codes?

type Activity[T any] struct {
	Context     string `json:"@context"`
	Type        string
	Integration string
	Actor       string
	Object      T
	Target      string
	Result      Result
}

type GenreralResult struct {
	Context     string `json:"@context"`
	Type        string
	Actor       string
	Object      map[string]any
	Target      string
	Result      Result
}

type Result struct {
	Type    string
	Code    string
	Summary string
}

func CreateResult(code int, summary string) Result {
	return Result{
		Type:    "httpResponseStatus",
		Code:    strconv.Itoa(code),
		Summary: summary,
	}
}
