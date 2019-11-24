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
		{
			"Nested fields",
			Person{
				"Bobo",
				Profile{33, "Shanghai"},
			},
			[]string{"Bobo", "Shanghai"},
		},
		{
			"Pointers to things",
			&Person{
				"Bobo",
				Profile{33, "Shanghai"},
			},
			[]string{"Bobo", "Shanghai"},
		},
		{
			"Slices",
			[]Profile{
				{33, "Shanghai"},
				{34, "Beijing"},
			},
			[]string{"Shanghai", "Beijing"},
		},
		{
			"Arrays",
			[2]Profile{
				{33, "Shanghai"},
				{34, "Beijing"},
			},
			[]string{"Shanghai", "Beijing"},
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

	t.Run("with maps", func(t *testing.T) {
		aMap := map[string]string{
			"Foo": "Bar",
			"Baz": "Boz",
		}

		var got []string
		walk(aMap, func(input string) {
			got = append(got, input)
		})

		assertContains(t, got, "Bar")
		assertContains(t, got, "Boz")
	})
}

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func assertContains(t *testing.T, haystack []string, needle string) {
	contains := false
	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}
	if !contains {
		t.Errorf("expected %+v to contain %q but it didn't", haystack, needle)
	}
}
