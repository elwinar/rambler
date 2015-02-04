package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

const (
	foo = `-- rambler up

create table foo (
	id integer
);

-- rambler down

drop table foo;
`
)

func Test_NewMigration_GlobError(t *testing.T) {
	_, err := NewMigration("unknown_directory", 0)

	if err == nil {
		t.Error("didn't failed on glob error")
	}
}

func Test_NewMigration_UnknownVersion(t *testing.T) {
	_, err := NewMigration("test", 0)

	if err == nil {
		t.Error("didn't failed on unknown version")
	}
}

func Test_NewMigration_AmbiguousVersion(t *testing.T) {
	_, err := NewMigration("test", 2)

	if err == nil {
		t.Error("didn't failed on ambiguous version")
	}
}

func Test_NewMigration_ParseDescription(t *testing.T) {
	type Case struct{
		Version uint64
		Description string
	}
	var cases = []Case{
		Case{3, "snake_case"},
		Case{4, "UpperCamelCase"},
		Case{5, "lowerCamelCase"},
		Case{6, "Title Case"},
		Case{7, "url-case"},
	}

	for n, c := range cases {
		m, err := NewMigration("test", c.Version)

		if err != nil {
			t.Error("unexpected error:", err)
		}

		if m.Description != c.Description {
			t.Errorf("uncorrectly parsed description for case %d: got \"%s\"", n, m.Description)
		}
	}
}

func Test_NewMigration_OK(t *testing.T) {
	m, err := NewMigration("test", 1)

	if err != nil {
		t.Error("unexpected error:", err)
	}

	if m.Version != 1 {
		t.Errorf("uncorrectly initialized migration version")
	}

	content, err := ioutil.ReadAll(m.Reader)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if string(content) != foo {
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
