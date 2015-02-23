package main

import (
	"reflect"
	"testing"
)

func Test_NewService(t *testing.T) {
	var cases = []struct {
		input Environment
		err   bool
	}{
		{
			input: Environment{
				Driver:    "mysql",
				Directory: "unkown",
			},
			err: true,
		},
		{
			input: Environment{
				Driver:    "unkown",
				Directory: "test",
			},
			err: true,
		},
		{
			input: Environment{
				Driver:    "mysql",
				Directory: "test",
			},
			err: false,
		},
	}

	for n, c := range cases {
		_, err := NewService(c.input)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
		}
	}
}

func Test_Service_ListAvailableMigrations_ParseFilenames(t *testing.T) {
	var cases = []struct {
		directory string
		output    []uint64
	}{
		{
			directory: "test/",
			output:    []uint64{1, 2, 3},
		},
		{
			directory: "test/empty/",
			output:    nil,
		},
		{
			directory: "test/not_a_directory",
			output:    nil,
		},
	}

	for n, c := range cases {
		s := CoreService{
			env: Environment{
				Directory: c.directory,
			},
		}
		versions := s.ListAvailableMigrations()

		if !reflect.DeepEqual(versions, c.output) {
			t.Error("case", n, "got unexpected outout:", versions)
		}
	}
}
