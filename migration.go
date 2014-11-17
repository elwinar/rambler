package rambler

import (
	"errors"
	"io"
	"os"
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
)

// NewMigration get a migration given its directory and version number
func NewMigration(directory string, version uint64) (*Migration, error) {
	if _, err := os.Stat(directory); err != nil {
		return nil, ErrUnknownDirectory
	}
	
	return nil, nil
}
