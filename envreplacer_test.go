package main

import (
	"os"
	"testing"
)

func TestFindEntries(t *testing.T) {
	{
		known := "INSERT INTO ${schema}.accounts (id, password) VALUES (2, '${pwd}')"
		found := findEntries(known)
		if len(found) != 2 {
			t.Error("Did not find entries")
			return
		}
	}

	{
		known := "INSERT INTO accounts (id, password) VALUES (2, 'iscleartext')"
		found := findEntries(known)
		if len(found) != 0 {
			t.Error("Did not find entries")
			return
		}
	}
}

func TestFindEnv(t *testing.T) {
	{
		os.Setenv("MYPASSWORD2", "alsocleartext")
		_, err := findEnvVal("${myPassword}")
		if err == nil {
			t.Error("should have failed")
		}
	}
	{
		os.Setenv("MYPASSWORD", "alsocleartext")
		_, err := findEnvVal("${myPassword}")
		if err != nil {
			t.Error("should have found")
		}
	}
}

func TestReplace(t *testing.T) {
	{
		os.Setenv("MYPASSWORD", "alsocleartext")
		os.Setenv("SCHEMA", "example")
		statement :=
			`INSERT INTO ${schema}.profile (id, password) VALUES (2, '${mypassword}');`
		rs, err := replace(statement)
		if err != nil {
			t.Error("Failed to replace ", err)
		}

		expected :=
			`INSERT INTO example.profile (id, password) VALUES (2, 'alsocleartext');`
		if expected != rs {
			t.Errorf("expected: '%s' does not match '%s'", expected, rs)
		}


		statement =
			`INSERT INTO ${schema}.profile (id, password) VALUES (2, '${mypassword}');
INSERT INTO ${schema}.profile (id, password) VALUES (2, '${mypassword}');`
		rs, err = replace(statement)
		if err != nil {
			t.Error("Failed to replace ", err)
		}
		expected2 := `INSERT INTO example.profile (id, password) VALUES (2, 'alsocleartext');
INSERT INTO example.profile (id, password) VALUES (2, 'alsocleartext');`
		if expected2 != rs {
			t.Errorf("expected: '%s' does not match '%s'", expected2, rs)
		}

		statement = "There is nothing to replace. Move along now."
		rs, err = replace(statement)
		if err != nil {
			t.Error("Failed to replace ", err)
		}
		if statement != rs {
			t.Errorf("Expected equal. '%s' and '%s' ", statement, rs)
		}
	}
}
