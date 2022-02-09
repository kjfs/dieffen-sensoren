package main

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf(fmt.Sprintf("Type is not http.Handler, it is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		log.Printf("Type of v %T", v)
		log.Printf("Type of h %T", h)
		log.Printf("Type of v %v", v)
		log.Printf("Type of h %v", h)
	default:
		t.Errorf(fmt.Sprintf("Type is not http.Handler, it is %T", v))
	}
}
