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
}

// NewMigration get a migration given its directory and version number
func NewMigration(directory string, version uint64) (*Migration, error) {
	matches, err := filepath.Glob(path.Join(directory, strconv.FormatUint(version, 10)+"_*.sql"))
	if err != nil {
		return nil, fmt.Errorf("directory %s unavailable: %s", directory, err.Error())
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no migration for version %d", version)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("ambiguous version %d", version)
	}

	m := &Migration{
		Name:        matches[0],
		Version:     version,
		Description: strings.Split(strings.SplitN(matches[0], "_", 2)[1], ".")[0],
	}

	return m, nil
}

// Scan retrieve all sections of the file with the given section marker.
func (m *Migration) Scan(section string) ([]string, error) {
	file, err := os.Open(m.Name)
	if err != nil {
		return nil, fmt.Errorf("file %s unavailable: %s", m.Name, err.Error())
	}
	return scan(file, section), nil
}

func scan(reader io.Reader, section string) []string {
	var scanner = bufio.NewScanner(reader)
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
