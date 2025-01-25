package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const prefix = `-- rambler`

// Migration represent a migration file, composed of up and down sections containing
// one or more statements each.
type Migration struct {
	Name string
	path string
}

// NewMigration generate a migration from the given file
func NewMigration(path string) (*Migration, error) {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("can't find migration %q: %w", path, err)
	}

	m := &Migration{
		Name: filepath.Base(path),
		path: path,
	}

	return m, nil
}

// Scan retrieve all sections of the file with the given section marker.
func (m Migration) scan(section string) ([]string, error) {
	file, err := os.Open(m.path)
	if err != nil {
		return nil, fmt.Errorf("opening migration %q: %w", m.path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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
	raw, err := m.scan(`down`)
	if err != nil {
		return nil, err
	}
	var stmt []string
	for i := len(raw) - 1; i >= 0; i-- {
		stmt = append(stmt, raw[i])
	}
	return stmt, nil
}
