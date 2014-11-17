package rambler

import (
	"errors"
	"io"
	"path"
	"path/filepath"
	"os"
	"strconv"
)

// Migration represent a migration file, composed of up and down sections containing
// one or more statements each.
type Migration struct {
	Version     uint64
	Description string
	reader      io.Reader
}

var (
	ErrUnknownDirectory = errors.New("unknown directory")
	ErrUnknownVersion = errors.New("unknwon version")
)

// NewMigration get a migration given its directory and version number
func NewMigration(directory string, version uint64) (migration *Migration, err error) {
	if _, err := os.Stat(directory); err != nil {
		return nil, ErrUnknownDirectory
	}
	
	matches, err := filepath.Glob(path.Join(directory, strconv.FormatUint(version, 10) + "_*.sql"))
	if err != nil {
		return nil, err
	}
	
	if len(matches) == 0 {
		return nil, ErrUnknownVersion
	}
	
	return nil, nil
}
