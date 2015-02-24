package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const prefix = `-- rambler`

// Migration represent a migration file, composed of up and down sections containing
// one or more statements each.
type Migration struct {
	Name        string
	Version     uint64
	Description string
	AppliedAt   *time.Time
	reader      io.Reader
}

// NewMigration get a migration given its directory and version number
func NewMigration(directory string, version uint64) (*Migration, error) {
	fi, err := os.Stat(directory)
	if err != nil {
		return nil, err
	}

	if !fi.Mode().IsDir() {
		return nil, fmt.Errorf("file %s isn't a directory", directory)
	}

	matches, _ := filepath.Glob(path.Join(directory, strconv.FormatUint(version, 10)+"_*.sql"))

	if len(matches) == 0 {
		return nil, fmt.Errorf("no migration for version %d", version)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("ambiguous version %d", version)
	}

	file, err := os.Open(matches[0])
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %s", matches[0], err.Error())
	}

	m := &Migration{
		Name:        matches[0],
		Version:     version,
		Description: strings.Split(strings.SplitN(matches[0], "_", 2)[1], ".")[0],
		reader:      file,
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

func (m Migration) Up() []string {
	return m.scan(`up`)
}

func (m Migration) Down() []string {
	raw := m.scan(`down`)
	var stmt []string
	for i := len(raw) - 1; i >= 0; i-- {
		stmt = append(stmt, raw[i])
	}
	return stmt
}
