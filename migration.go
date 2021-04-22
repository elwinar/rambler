package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const prefix = `-- rambler`

// Migration represent a migration file, composed of up and down sections containing
// one or more statements each.
type Migration struct {
	Path string
	Name string
}

// NewMigration generate a migration from the given file
func NewMigration(path string) (Migration, error) {
	_, err := os.Stat(path)
	if err != nil {
		return Migration{}, fmt.Errorf("unable to find file %s: %s", path, err.Error())
	}

	return Migration{
		Path: path,
		Name: filepath.Base(path),
	}, nil
}

// Scan retrieve all sections of the file with the given section marker.
func (m Migration) scan(section string) ([]string, error) {
	file, err := os.Open(m.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %s", m.Name, err.Error())
	}
	defer file.Close()

	var scanner = bufio.NewScanner(file)
	var statements []string
	var buffer string

	recording := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, prefix) {
			if len(strings.TrimSpace(buffer)) != 0 {
				statements = append(statements, strings.TrimSpace(buffer))
			}

			buffer = ""
			cmd := strings.TrimSpace(line[len(prefix):])

			switch cmd {
			case section:
				recording = true
			default:
				recording = false
			}

			continue
		}

		if recording {
			buffer = buffer + "\n" + line
		}
	}

	if len(strings.TrimSpace(buffer)) != 0 {
		statements = append(statements, strings.TrimSpace(buffer))
	}

	return statements, nil
}

// Up return the up statements of the migration
func (m Migration) Up() ([]string, error) {
	return m.scan(`up`)
}

// Down return the down statements of the migration
func (m Migration) Down() ([]string, error) {
	stmts, err := m.scan(`down`)
	for i, j := 0, len(stmts)-1; i < j; i, j = i+1, j-1 {
		stmts[i], stmts[j] = stmts[j], stmts[i]
	}
	return stmts, err
}
