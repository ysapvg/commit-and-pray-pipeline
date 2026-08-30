package main

import (
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	expected := "Hello, DevOps! version=dev\n"

	if rec.Body.String() != expected {
		t.Errorf("got %q, want %q", rec.Body.String(), expected)
	}
}
