package rambler

import (
	"bufio"
	"errors"
	"io"
	"path"
	"os"
	"strconv"
	"strings"
)

// Migration represent a migration file, composed of up and down sections containing
// one or more statements each.
type Migration struct {
	Version     uint64
	Description string
	reader      io.ReadSeeker
}

const (
	prefix = "-- rambler"
)

var (
	ErrUnknownDirectory = errors.New("unknown directory")
	ErrUnknownVersion = errors.New("unknwon version")
	ErrAmbiguousVersion = errors.New("ambiguous version")
)

// NewMigration get a migration given its directory and version number
func NewMigration(directory string, version uint64, glober Glober, opener Opener) (*Migration, error) {
	if _, err := os.Stat(directory); err != nil {
		return nil, ErrUnknownDirectory
	}
	
	matches, err := glober.Glob(path.Join(directory, strconv.FormatUint(version, 10) + "_*.sql"))
	if err != nil {
		return nil, err
	}
	
	if len(matches) == 0 {
		return nil, ErrUnknownVersion
	}
	
	if len(matches) > 1 {
		return nil, ErrAmbiguousVersion
	}
	
	reader, err := opener.Open(matches[0])
	if err != nil {
		return nil, err
	}
	
	m := &Migration{
		Version: version,
		Description: strings.Split(strings.SplitN(matches[0], "_", 2)[1], ".")[0],
		reader: reader,
	}
	
	return m, nil
}

// Scan retrieve all sections of the file with the given section marker.
func (m *Migration) Scan(section string) ([]string, error) {
	m.reader.Seek(0,0)
	
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
	
	return statements, nil
}
