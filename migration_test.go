package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_NewMigration(t *testing.T) {
	os.Chmod("test/4_unreadable.sql", 0000)
	
	var cases = []struct {
		directory string
		version   uint64
		err       bool
		output    *Migration
	}{
		{
			directory: "unknown_directory",
			version:   0,
			err:       true,
			output:    nil,
		},
		{
			directory: "test/not_a_directory",
			version:   0,
			err:       true,
			output:    nil,
		},
		{
			directory: "test",
			version:   0,
			err:       true,
			output:    nil,
		},
		{
			directory: "test",
			version:   1,
			err:       false,
			output: &Migration{
				Name:        "test/1_foo.sql",
				Version:     1,
				Description: "foo",
				AppliedAt:   nil,
			},
		},
		{
			directory: "test",
			version:   2,
			err:       true,
			output:    nil,
		},
		{
			directory: "test",
			version:   4,
			err:       true,
			output:    nil,
		},
	}

	for n, c := range cases {
		migration, err := NewMigration(c.directory, c.version)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if migration != nil {
			migration.reader = nil
		}

		if !reflect.DeepEqual(migration, c.output) {
			t.Error("case", n, "got unexpected output:", migration)
		}
	}
	
	os.Chmod("test/4_unreadable.sql", 0664)
}

func Test_Migration_Up(t *testing.T) {
	var cases = []struct {
		reader  *strings.Reader
		output  []string
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
			output:  []string{"first", "second", "fourth"},
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


func Test_Migration_Down(t *testing.T) {
	var cases = []struct {
		reader  *strings.Reader
		output  []string
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
			output:  []string{"third"},
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
			output:  []string{"third", "second"},
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
