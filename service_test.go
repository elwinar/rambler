package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestNewService(t *testing.T) {
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
				Driver:    "mysql",
				Directory: "test/not_a_directory",
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

func TestServiceAvailable(t *testing.T) {
	var cases = []struct {
		directory  string
		migrations []*Migration
		err        bool
	}{
		{
			directory:  "test/empty",
			migrations: nil,
			err:        false,
		},
		{
			directory: "test/one",
			migrations: []*Migration{
				&Migration{
					Name:   "1_one.sql",
					reader: nil,
				},
			},
			err: false,
		},
		{
			directory: "test/two",
			migrations: []*Migration{
				&Migration{
					Name:   "1_one.sql",
					reader: nil,
				},
				&Migration{
					Name:   "2_two.sql",
					reader: nil,
				},
			},
			err: false,
		},
		{
			directory: "test/others",
			migrations: []*Migration{
				&Migration{
					Name:   "1_one.sql",
					reader: nil,
				},
				&Migration{
					Name:   "2_two.sql",
					reader: nil,
				},
			},
			err: false,
		},
	}

	for n, c := range cases {
		service := &Service{
			env: Environment{
				Directory: c.directory,
			},
		}

		migrations, err := service.Available()
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
		}
		
		for _, m := range migrations {
			m.reader = nil
		}

		if !reflect.DeepEqual(migrations, c.migrations) {
			t.Error("case", n, "got unexpected migrations:", migrations)
		}
	}
}

func TestServiceApplied(t *testing.T) {
	var cases = []struct{
		directory string
		table []string
		fail error
		migrations []*Migration
		err bool
	}{
		{
			directory: "test/one",
			table: []string{
				"1_one.sql",
			},
			fail: nil,
			migrations: []*Migration{
				&Migration{
					Name:   "1_one.sql",
					reader: nil,
				},
			},
			err: false,
		},
		{
			directory: "test/one",
			table: []string{},
			fail: errors.New("error"),
			migrations: nil,
			err: true,
		},
		{
			directory: "test/one",
			table: []string{
				"1_one.sql",
				"2_two.sql",
			},
			fail: nil,
			migrations: nil,
			err: true,
		},
	}

	for n, c := range cases {
		service := &Service{
			env: Environment{
				Directory: c.directory,
			},
			conn: MockConn{
				getApplied: func() ([]string, error) {
					return c.table, c.fail
				},
			},
		}

		migrations, err := service.Applied()
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
		}
		
		for _, m := range migrations {
			m.reader = nil
		}

		if !reflect.DeepEqual(migrations, c.migrations) {
			t.Error("case", n, "got unexpected migrations:", migrations)
		}
	}
}

func TestServiceApply(t *testing.T) {
	var cases = []struct{
		migration *Migration
		executeFail error
		addAppliedFail error
		err bool
		executed []string
		applied []string
	}{
		{
			migration: &Migration{
				Name: "1_one.sql",
				reader: strings.NewReader(`-- rambler up
first
-- rambler up
second
-- rambler down
third
-- rambler up
fourth
`),
			},
			executeFail: nil,
			addAppliedFail: nil,
			err: false,
			executed: []string{
				"first",
				"second",
				"fourth",
			},
			applied: []string{
				"1_one.sql",
			},
		},
		{
			migration: &Migration{
				Name: "1_one.sql",
				reader: strings.NewReader(`-- rambler up
first
-- rambler up
second
-- rambler down
third
-- rambler up
fourth
`),
			},
			executeFail: errors.New("error"),
			addAppliedFail: nil,
			err: true,
			executed: []string{
				"first",
			},
			applied: nil,
		},
		{
			migration: &Migration{
				Name: "1_one.sql",
				reader: strings.NewReader(`-- rambler up
first
-- rambler up
second
-- rambler down
third
-- rambler up
fourth
`),
			},
			executeFail: nil,
			addAppliedFail: errors.New("error"),
			err: true,
			executed: []string{
				"first",
				"second",
				"fourth",
			},
			applied: []string{
				"1_one.sql",
			},
		},
	}

	for n, c := range cases {
		var executed, applied []string
		service := &Service{
			conn: MockConn{
				execute: func(statement string) error {
					executed = append(executed, statement)
					return c.executeFail
				},
				addApplied: func(migration string) error {
					applied = append(applied, migration)
					return c.addAppliedFail
				},
			},
		}

		err := service.Apply(c.migration)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
		}

		if !reflect.DeepEqual(executed, c.executed) {
			t.Error("case", n, "got unexpected statements:", executed)
		}

		if !reflect.DeepEqual(applied, c.applied) {
			t.Error("case", n, "got unexpected applied:", applied)
		}
	}
}
