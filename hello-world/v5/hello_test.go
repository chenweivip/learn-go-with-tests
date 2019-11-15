package main

import "testing"

func TestHello(t *testing.T) {

	assertCorrectMessage := func(t *testing.T, got, expected string) {
		t.Helper()
		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	}

	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Bobo")
		expected := "Hello, Bobo"
		assertCorrectMessage(t, got, expected)
	})

	t.Run("empty string defaults to 'world'", func(t *testing.T) {
		got := Hello("")
		expected := "Hello, World1"
		assertCorrectMessage(t, got, expected)
	})

}
