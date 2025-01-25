package main

import (
	"reflect"
	"testing"
)

func TestNewMigration(t *testing.T) {
	cases := []struct {
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
				Name: "1_foo.sql",
				path: "testdata/1_foo.sql",
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

		if !reflect.DeepEqual(migration, c.migration) {
			t.Error("case", n, "got unexpected output:", migration)
		}
	}
}

func TestMigrationUp(t *testing.T) {
	cases := []struct {
		path   string
		output []string
	}{
		{
			path:   "testdata/migrations/01_up.sql",
			output: []string{"first", "second", "fourth"},
		},
	}

	for n, c := range cases {
		migration := &Migration{
			path: c.path,
		}

		statements, err := migration.Up()
		if err != nil {
			t.Error("case", n, "got unexpected error:", err)
		}
		if !reflect.DeepEqual(statements, c.output) {
			t.Error("case", n, "got unexpected output:", statements)
		}
	}
}

func TestMigrationDown(t *testing.T) {
	cases := []struct {
		path   string
		output []string
	}{
		{
			path:   "testdata/migrations/02_down1.sql",
			output: []string{"third"},
		},
		{
			path:   "testdata/migrations/02_down2.sql",
			output: []string{"third", "second"},
		},
	}

	for n, c := range cases {
		migration := &Migration{
			path: c.path,
		}

		statements, err := migration.Down()
		if err != nil {
			t.Error("case", n, "got unexpected error:", err)
		}
		if !reflect.DeepEqual(statements, c.output) {
			t.Error("case", n, "got unexpected output:", statements)
		}
	}
}
