package main

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Bobo")

	got := buffer.String()
	expected := "Hello, Bobo"

	if got != expected {
		t.Errorf("got %q expected %q", got, expected)
	}
}
