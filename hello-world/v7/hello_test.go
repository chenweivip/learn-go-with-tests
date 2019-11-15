package main

import "testing"

func TestHello(t *testing.T) {

	assertCorrectMessage := func(got, expected string) {
		t.Helper()
		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	}

	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Bobo", "")
		expected := "Hello, Bobo"
		assertCorrectMessage(got, expected)
	})

	t.Run("say hello world when an empty string is supplied", func(t *testing.T) {
		got := Hello("", "")
		expected := "Hello, World"
		assertCorrectMessage(got, expected)
	})

	t.Run("say hello in Chinese", func(t *testing.T) {
		got := Hello("波波", chinese)
		expected := "你好, 波波"
		assertCorrectMessage(got, expected)
	})

	t.Run("say hello in French", func(t *testing.T) {
		got := Hello("Lauren", french)
		expected := "Bonjour, Lauren"
		assertCorrectMessage(got, expected)
	})

}
