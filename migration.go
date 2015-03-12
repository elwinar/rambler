package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const prefix = `-- rambler`

// Migration represent a migration file, composed of up and down sections containing
// one or more statements each.
type Migration struct {
	Name   string
	reader io.Reader
}

// NewMigration generate a migration from the given file
func NewMigration(path string) (*Migration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %s", path, err.Error())
	}

	m := &Migration{
		Name:      filepath.Base(path),
		reader:    file,
	}

	return m, nil
}

// Scan retrieve all sections of the file with the given section marker.
func (m Migration) scan(section string) []string {
	var scanner = bufio.NewScanner(m.reader)
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

	return statements
}

// Return the up statements of the migration
func (m Migration) Up() []string {
	return m.scan(`up`)
}

// Return the down statements of the migration
func (m Migration) Down() []string {
	raw := m.scan(`down`)
	var stmt []string
	for i := len(raw) - 1; i >= 0; i-- {
		stmt = append(stmt, raw[i])
	}
	return stmt
}
