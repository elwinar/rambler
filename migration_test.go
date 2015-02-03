package main

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_NewMigration_GlobError(t *testing.T) {
	_, err := newMigration("test", 0, func(_ string) ([]string, error) {
		return nil, errors.New("glob error")
	}, func(_ string) (io.ReadSeeker, error) {
		return strings.NewReader("ok"), nil
	})

	if err == nil {
		t.Error("didn't failed on glob error")
	}

	if err.Error() != "directory test unavailable: glob error" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_NewMigration_UnknownVersion(t *testing.T) {
	_, err := newMigration("", 0, func(_ string) ([]string, error) {
		return []string{}, nil
	}, func(_ string) (io.ReadSeeker, error) {
		return strings.NewReader("ok"), nil
	})

	if err == nil {
		t.Error("didn't failed on unknown version")
	}

	if err.Error() != "no migration for version 0" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_NewMigration_AmbiguousVersion(t *testing.T) {
	_, err := newMigration("", 0, func(_ string) ([]string, error) {
		return []string{"0_first.sql", "0_second.sql"}, nil
	}, func(_ string) (io.ReadSeeker, error) {
		return strings.NewReader("ok"), nil
	})

	if err == nil {
		t.Error("didn't failed on ambiguous version")
	}

	if err.Error() != "ambiguous version 0" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_NewMigration_OpenError(t *testing.T) {
	_, err := newMigration("", 0, func(_ string) ([]string, error) {
		return []string{"0_first.sql"}, nil
	}, func(_ string) (io.ReadSeeker, error) {
		return nil, errors.New("open error")
	})

	if err == nil {
		t.Error("didn't failed on open error")
	}

	if err.Error() != "file 0_first.sql unavailable: open error" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_NewMigration_ParseDescription(t *testing.T) {
	var cases = []string{
		"snake_case",
		"UpperCamelCase",
		"lowerCamelCase",
		"Title Case",
	}

	for n, c := range cases {
		m, err := newMigration("", 0, func(_ string) ([]string, error) {
			return []string{"0_" + c + ".sql"}, nil
		}, func(_ string) (io.ReadSeeker, error) {
			return nil, nil
		})

		if err != nil {
			t.Error("unexpected error:", err)
		}

		if m.Description != c {
			t.Errorf("uncorrectly parsed description for case %d: got \"%s\"", n, m.Description)
		}
	}
}

func Test_NewMigration_OK(t *testing.T) {
	m, err := newMigration("", 0, func(_ string) ([]string, error) {
		return []string{"0_description.sql"}, nil
	}, func(_ string) (io.ReadSeeker, error) {
		return strings.NewReader("migration"), nil
	})

	if err != nil {
		t.Error("unexpected error:", err)
	}

	if m.Version != 0 {
		t.Errorf("uncorrectly initialized migration version")
	}

	content, err := ioutil.ReadAll(m.Reader)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if string(content) != "migration" {
		t.Errorf("uncorrectly initialized migration reader")
	}
}

func Test_Scan_ParseStatements(t *testing.T) {
	type Case struct {
		Content    string
		Section    string
		Statements []string
	}
	var cases = []Case{
		Case{
			Content: `-- rambler up
one
-- rambler up
two
`,
			Section:    "up",
			Statements: []string{"one", "two"},
		},
		Case{
			Content: `-- rambler up
one
-- rambler down
two
-- rambler up
three
`,
			Section:    "up",
			Statements: []string{"one", "three"},
		},
		Case{
			Content: `-- rambler up
one
-- rambler down
two
-- rambler up
three
`,
			Section:    "down",
			Statements: []string{"two"},
		},
	}

	for n, c := range cases {
		m := &Migration{
			Reader: strings.NewReader(c.Content),
		}
		statements := m.Scan(c.Section)

		if len(statements) != len(c.Statements) {
			t.Error("found incorrect number of statements for case %d: got", n, len(statements))
			continue
		}

		for i := 0; i < len(statements); i++ {
			if statements[i] != c.Statements[i] {
				t.Error("didn't parsed correctly statement %d of case %d", i, n)
			}
		}
	}
}
