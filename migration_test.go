package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewMigration(t *testing.T) {
	var cases = []struct {
		path      string
		migration *Migration
		err       bool
	}{
		{
			path:      "testdata/0_unknown.sql",
			migration: nil,
			err:       true,
		},
		{
			path: "testdata/1_foo.sql",
			migration: &Migration{
				Name:   "1_foo.sql",
				reader: nil,
			},
			err: false,
		},
	}

	for n, c := range cases {
		migration, err := NewMigration(c.path)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if migration != nil {
			migration.reader = nil
		}

		if !reflect.DeepEqual(migration, c.migration) {
			t.Error("case", n, "got unexpected output:", migration)
		}
	}
}

func TestMigrationUp(t *testing.T) {
	var cases = []struct {
		reader *strings.Reader
		output []string
	}{
		{
			reader: strings.NewReader(`-- rambler up
first
-- rambler up
second
-- rambler down
third
-- rambler up
fourth
`),
			output: []string{"first", "second", "fourth"},
		},
	}

	for n, c := range cases {
		migration := &Migration{
			reader: c.reader,
		}

		statements := migration.Up()
		if !reflect.DeepEqual(statements, c.output) {
			t.Error("case", n, "got unexpected output:", statements)
		}
	}
}

func TestMigrationDown(t *testing.T) {
	var cases = []struct {
		reader *strings.Reader
		output []string
	}{
		{
			reader: strings.NewReader(`-- rambler up
first
-- rambler up
second
-- rambler down
third
-- rambler up
fourth
`),
			output: []string{"third"},
		},
		{
			reader: strings.NewReader(`-- rambler up
first
-- rambler down
second
-- rambler down
third
-- rambler up
fourth
`),
			output: []string{"third", "second"},
		},
	}

	for n, c := range cases {
		migration := &Migration{
			reader: c.reader,
		}

		statements := migration.Down()
		if !reflect.DeepEqual(statements, c.output) {
			t.Error("case", n, "got unexpected output:", statements)
		}
	}
}
