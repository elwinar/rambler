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
	path   string
	file   *os.File
}

func (m Migration) Close() error {
	if m.file != nil {
		return m.file.Close()
	}
	return nil
}

func (m Migration) r() io.Reader {
	if m.reader == nil {
		file, _ := os.Open(m.path) // Assumes the path was verified via NewMigration.
		m.reader = file
	}

	return m.reader
}

// NewMigration generate a migration from the given file
func NewMigration(path string) (*Migration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %s", path, err.Error())
	}
	if err = file.Close(); err != nil {
		return nil, fmt.Errorf("unable to close file %s: %s", path, err.Error())
	}

	m := &Migration{
		Name: filepath.Base(path),
		path: path,
	}

	return m, nil
}

// Scan retrieve all sections of the file with the given section marker.
func (m Migration) scan(section string) []string {
	defer m.Close()
	var scanner = bufio.NewScanner(m.r())
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

// Up return the up statements of the migration
func (m Migration) Up() []string {
	return m.scan(`up`)
}

// Down return the down statements of the migration
func (m Migration) Down() []string {
	raw := m.scan(`down`)
	var stmt []string
	for i := len(raw) - 1; i >= 0; i-- {
		stmt = append(stmt, raw[i])
	}
	return stmt
}
