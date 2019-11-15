package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Bobo")
	expected := "Hello, Bobo"

	if got != expected {
		t.Errorf("got %q expected %q", got, expected)
	}
}
