package migration

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// The various errors returned by the package
var (
	ErrUnknownDirectory   = errors.New("unknown directory")
	ErrUnknownVersion     = errors.New("unknwon version")
	ErrUnknownDriver      = errors.New("unknwon driver")
	ErrUnknownEnvironment = errors.New("unknwon environment")
	ErrAmbiguousVersion   = errors.New("ambiguous version")
)


// Migration represent a migration file, composed of up and down sections containing
// one or more statements each.
type Migration struct {
	Version     uint64
	Description string
	AppliedAt   *time.Time
	reader      io.ReadSeeker
}

const (
	prefix = "-- rambler"
)

// NewMigration get a migration given its directory and version number
func NewMigration(directory string, version uint64) (*Migration, error) {
	return newMigration(directory, version, filepath.Glob, func(path string) (io.ReadSeeker, error) {
		return os.Open(path)
	})
}

type glober func(string) ([]string, error)
type opener func(string) (io.ReadSeeker, error)

func newMigration(directory string, version uint64, glob glober, open opener) (*Migration, error) {
	matches, err := glob(path.Join(directory, strconv.FormatUint(version, 10)+"_*.sql"))
	if err != nil {
		return nil, ErrUnknownDirectory
	}

	if len(matches) == 0 {
		return nil, ErrUnknownVersion
	}

	if len(matches) > 1 {
		return nil, ErrAmbiguousVersion
	}

	reader, err := open(matches[0])
	if err != nil {
		return nil, err
	}

	m := &Migration{
		Version:     version,
		Description: strings.Split(strings.SplitN(matches[0], "_", 2)[1], ".")[0],
		reader:      reader,
	}

	return m, nil
}

// Scan retrieve all sections of the file with the given section marker.
func (m *Migration) Scan(section string) []string {
	m.reader.Seek(0, 0)

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
