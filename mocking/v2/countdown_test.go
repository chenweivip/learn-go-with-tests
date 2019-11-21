package main

import (
	"bytes"
	"testing"
)

func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}

	Countdown(buffer)

	got := buffer.String()
	expected := `3
2
1
Go!`

	if got != expected {
		t.Errorf("got %q expected %q", got, expected)
	}
}
