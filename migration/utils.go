package migration

import (
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration/driver"
	"io"
	"os"
)

const (
	prefix = "-- rambler"
)

var (
	errUnknownVersion       = "no migration for version %d"
	errAmbiguousVersion     = "ambiguous version %d"
	errUnavailableFile      = "file %s unavailable: %s"
	errDriverError          = "unable to initialize driver %s: %s"
	errUnavailableDirectory = "directory %s unavailable: %s"
)

type glober func(string) ([]string, error)
type opener func(string) (io.ReadSeeker, error)
type stater func(string) (os.FileInfo, error)
type connConstructor func(configuration.Environment) (driver.Conn, error)
