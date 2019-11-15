package main

import "testing"

func TestHello(t *testing.T) {

	assertCorrectMessage := func(t *testing.T, got, expected string) {
		t.Helper()
		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	}

	t.Run("to a person", func(t *testing.T) {
		got := Hello("Bobo", "")
		expected := "Hello, Bobo"
		assertCorrectMessage(t, got, expected)
	})

	t.Run("empty string", func(t *testing.T) {
		got := Hello("", "")
		expected := "Hello, World"
		assertCorrectMessage(t, got, expected)
	})

	t.Run("in Chinese", func(t *testing.T) {
		got := Hello("波波", chinese)
		expected := "你好, 波波"
		assertCorrectMessage(t, got, expected)
	})

	t.Run("in French", func(t *testing.T) {
		got := Hello("Lauren", french)
		expected := "Bonjour, Lauren"
		assertCorrectMessage(t, got, expected)
	})

}
