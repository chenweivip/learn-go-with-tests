package main

import "testing"

func TestPerimeter(t *testing.T) {
	got := Perimeter(10.0, 10.0)
	expected := 40.0

	if got != expected {
		t.Errorf("got %.2f expected %.2f", got, expected)
	}
}

func TestArea(t *testing.T) {
	got := Area(12.0, 6.0)
	expected := 72.0

	if got != expected {
		t.Errorf("got %.2f expected %.2f", got, expected)
	}
}
