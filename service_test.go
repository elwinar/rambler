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

func Test_Service_ListAvailableMigrations(t *testing.T) {
	var cases = []struct {
		directory string
		err       bool
		output    []uint64
	}{
		{
			directory: "test/",
			err:       false,
			output:    []uint64{1, 2, 3},
		},
		{
			directory: "test/empty/",
			err:       false,
			output:    nil,
		},
		{
			directory: "test/not_a_directory",
			err:       true,
			output:    nil,
		},
		{
			directory: "test/doesnt_exists",
			err:       true,
			output:    nil,
		},
	}

	for n, c := range cases {
		s := &Service{
			env: Environment{
				Directory: c.directory,
			},
		}

		versions, err := s.ListAvailableMigrations()
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
		}

		if !reflect.DeepEqual(versions, c.output) {
			t.Error("case", n, "got unexpected outout:", versions)
		}
	}
}
