package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurve(t *testing.T) {
	var myH myHandler

	h := NoSurve(&myH)

	switch v := h.(type) {
	case http.Handler:

	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:

	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}
