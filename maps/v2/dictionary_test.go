package main

import (
	"testing"
)

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		got, _ := dictionary.Search("test")
		expected := "this is just a test"

		assertStrings(t, got, expected)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, got := dictionary.Search("unknown")

		assertError(t, got, ErrNotFound)
	})
}

func assertStrings(t *testing.T, got, expected string) {
	t.Helper()

	if got != expected {
		t.Errorf("got %q expected %q", got, expected)
	}
}

func assertError(t *testing.T, got, expected error) {
	t.Helper()

	if got != expected {
		t.Errorf("got error %q expected %q", got, expected)
	}
}
