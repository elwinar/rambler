package main

import (
	"reflect"
	"testing"
)

func Test_NewMigration(t *testing.T) {
	var cases = []struct{
		directory string
		version uint64
		err bool
		output *Migration
	}{
		{
			directory: "unknown_directory",
			version: 0,
			err: true,
			output: nil,
		},
		{
			directory: "test",
			version: 0,
			err: true,
			output: nil,
		},
		{
			directory: "test",
			version: 1,
			err: false,
			output: &Migration{
				Name: "test/1_foo.sql",
				Version: 1,
				Description: "foo",
				AppliedAt: nil,
			},
		},
		{
			directory: "test",
			version: 2,
			err: true,
			output: nil,
		},
	}
	
	for n, c := range cases {
		migration, err := NewMigration(c.directory, c.version)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}
		
		if !reflect.DeepEqual(migration, c.output) {
			t.Error("case", n, "got unexpected output:", migration)
		}
	}
}
