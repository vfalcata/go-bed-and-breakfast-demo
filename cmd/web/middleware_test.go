package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH) // should accept a handler, and return one as well

	// test that the return type is correct, http.Handler
	switch v := h.(type) { // storing in variable "v" whatever type "h" is
	case http.Handler: // if the case  is http.Handler do nothing since it is the correct type
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but type is %T", v)) // %T is the type

	}

}

func TestSesssionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH) // should accept a handler, and return one as well

	// test that the return type is correct, http.Handler
	switch v := h.(type) { // storing in variable "v" whatever type "h" is
	case http.Handler: // if the case  is http.Handler do nothing since it is the correct type
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but type is %T", v)) // %T is the type

	}

}
