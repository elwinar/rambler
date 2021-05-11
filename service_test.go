package main

import (
	"errors"
	"reflect"
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
				Directory: "testdata/not_a_directory",
			},
			err: true,
		},
		{
			input: Environment{
				Driver:    "unkown",
				Directory: "testdata",
			},
			err: true,
		},
		{
			input: Environment{
				Driver:    "mysql",
				Directory: "testdata",
			},
		},
	}

	for n, c := range cases {
		_, err := NewService(c.input, false)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
		}
	}
}

func TestServiceAvailable(t *testing.T) {
	var cases = []struct {
		directory  string
		migrations []Migration
		err        bool
	}{
		{
			directory: "testdata/empty",
		},
		{
			directory: "testdata/one",
			migrations: []Migration{
				{
					Path: "testdata/one/1_one.sql",
					Name: "1_one.sql",
				},
			},
		},
		{
			directory: "testdata/two",
			migrations: []Migration{
				{
					Path: "testdata/two/1_one.sql",
					Name: "1_one.sql",
				},
				{
					Path: "testdata/two/2_two.sql",
					Name: "2_two.sql",
				},
			},
		},
		{
			directory: "testdata/others",
			migrations: []Migration{
				{
					Path: "testdata/others/1_one.sql",
					Name: "1_one.sql",
				},
				{
					Path: "testdata/others/2_two.sql",
					Name: "2_two.sql",
				},
			},
		},
		{
			directory: "testdata/unreachable",
			err:       true,
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

		if !reflect.DeepEqual(migrations, c.migrations) {
			t.Error("case", n, "got unexpected migrations:", migrations)
		}
	}
}

func TestServiceApplied(t *testing.T) {
	var cases = []struct {
		directory  string
		table      []string
		fail       error
		migrations []Migration
		err        bool
	}{
		{
			directory: "testdata/one",
			table: []string{
				"1_one.sql",
			},
			migrations: []Migration{
				{
					Path: "testdata/one/1_one.sql",
					Name: "1_one.sql",
				},
			},
		},
		{
			directory: "testdata/one",
			table:     []string{},
			fail:      errors.New("error"),
			err:       true,
		},
		{
			directory: "testdata/one",
			table: []string{
				"1_one.sql",
				"2_two.sql",
			},
			err: true,
		},
		{
			directory: "testdata/two",
			table: []string{
				"1_one.sql",
				"2_two.sql",
			},
			migrations: []Migration{
				{
					Path: "testdata/two/1_one.sql",
					Name: "1_one.sql",
				},
				{
					Path: "testdata/two/2_two.sql",
					Name: "2_two.sql",
				},
			},
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

		if !reflect.DeepEqual(migrations, c.migrations) {
			t.Error("case", n, "got unexpected migrations:", migrations)
		}
	}
}

func TestServiceApply(t *testing.T) {
	var cases = []struct {
		migration      Migration
		executeFail    error
		addAppliedFail error
		save           bool
		err            bool
		executed       []string
		applied        []string
	}{
		{
			migration: Migration{
				Path: "testdata/dummy.sql",
				Name: "dummy.sql",
			},
			save: true,
			executed: []string{
				"first",
				"second",
				"fourth",
			},
			applied: []string{"dummy.sql"},
		},
		{
			migration: Migration{
				Path: "testdata/dummy.sql",
				Name: "dummy.sql",
			},
			executeFail: errors.New("error"),
			save:        true,
			err:         true,
			executed: []string{
				"first",
			},
		},
		{
			migration: Migration{
				Path: "testdata/dummy.sql",
				Name: "dummy.sql",
			},
			addAppliedFail: errors.New("error"),
			save:           true,
			err:            true,
			executed: []string{
				"first",
				"second",
				"fourth",
			},
			applied: []string{"dummy.sql"},
		},
		{
			err: true,
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

		err := service.Apply(c.migration, c.save)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
		}

		if !reflect.DeepEqual(executed, c.executed) {
			t.Errorf("case %d got unexpected statements: wanted %+v, got %+v", n, c.executed, executed)
		}

		if !reflect.DeepEqual(applied, c.applied) {
			t.Errorf("case %d got unexpected applied: wanted %+v, got %+v", n, c.applied, applied)
		}
	}
}

func TestServiceReverse(t *testing.T) {
	var cases = []struct {
		migration         Migration
		executeFail       error
		removeAppliedFail error
		save              bool
		err               bool
		executed          []string
		reversed          []string
	}{
		{
			migration: Migration{
				Path: "testdata/dummy.sql",
				Name: "dummy.sql",
			},
			save: true,
			executed: []string{
				"fifth",
				"third",
			},
			reversed: []string{"dummy.sql"},
		},
		{
			migration: Migration{
				Path: "testdata/dummy.sql",
				Name: "dummy.sql",
			},
			save:        true,
			executeFail: errors.New("error"),
			err:         true,
			executed: []string{
				"fifth",
			},
		},
		{
			migration: Migration{
				Path: "testdata/dummy.sql",
				Name: "dummy.sql",
			},
			removeAppliedFail: errors.New("error"),
			save:              true,
			err:               true,
			executed: []string{
				"fifth",
				"third",
			},
			reversed: []string{"dummy.sql"},
		},
		{
			err: true,
		},
	}

	for n, c := range cases {
		var executed, reversed []string
		service := &Service{
			conn: MockConn{
				execute: func(statement string) error {
					executed = append(executed, statement)
					return c.executeFail
				},
				removeApplied: func(migration string) error {
					reversed = append(reversed, migration)
					return c.removeAppliedFail
				},
			},
		}

		err := service.Reverse(c.migration, c.save)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
		}

		if !reflect.DeepEqual(executed, c.executed) {
			t.Errorf("case %d got unexpected statements: wanted %+v, got %+v", n, c.executed, executed)
		}

		if !reflect.DeepEqual(reversed, c.reversed) {
			t.Errorf("case %d got unexpected reversed: wanted %+v, got %+v", n, c.reversed, reversed)
		}
	}
}
