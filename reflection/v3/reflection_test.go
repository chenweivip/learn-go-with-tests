package main

import (
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {

	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"Struct with one string field",
			struct{ Name string }{"Bobo"},
			[]string{"Bobo"},
		},
		{
			"Struct with two string fields",
			struct {
				Name string
				City string
			}{"Bobo", "Shanghai"},
			[]string{"Bobo", "Shanghai"},
		},
		{
			"Struct with non string field",
			struct {
				Name string
				Age  int
			}{"Bobo", 33},
			[]string{"Bobo"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, expect %v", got, test.ExpectedCalls)
			}
		})
	}
}
