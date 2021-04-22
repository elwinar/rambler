package main

import (
	"reflect"
	"testing"
)

func TestNewMigration(t *testing.T) {
	var cases = []struct {
		path      string
		migration Migration
		err       bool
	}{
		{
			path: "testdata/0_unknown.sql",
			err:  true,
		},
		{
			path: "testdata/1_foo.sql",
			migration: Migration{
				Path: "testdata/1_foo.sql",
				Name: "1_foo.sql",
			},
		},
	}

	for n, c := range cases {
		migration, err := NewMigration(c.path)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if !reflect.DeepEqual(migration, c.migration) {
			t.Error("case", n, "got unexpected output:", migration)
		}
	}
}

func TestMigrationUp(t *testing.T) {
	var cases = []struct {
		migration Migration
		output    []string
		err       bool
	}{
		{
			migration: Migration{
				Path: "testdata/dummy.sql",
				Name: "dummy.sql",
			},
			output: []string{"first", "second", "fourth"},
		},
	}

	for n, c := range cases {
		statements, err := c.migration.Up()
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if !reflect.DeepEqual(statements, c.output) {
			t.Error("case", n, "got unexpected output:", statements)
		}
	}
}

func TestMigrationDown(t *testing.T) {
	var cases = []struct {
		migration Migration
		output    []string
		err       bool
	}{
		{
			migration: Migration{
				Path: "testdata/dummy.sql",
				Name: "dummy.sql",
			},
			output: []string{"fifth", "third"},
		},
	}

	for n, c := range cases {
		statements, err := c.migration.Down()
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if !reflect.DeepEqual(statements, c.output) {
			t.Error("case", n, "got unexpected output:", statements)
		}
	}
}
