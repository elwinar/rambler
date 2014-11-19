package migration

import (
	"errors"
)

// The various errors returned by the package
var (
	ErrUnknownDirectory = errors.New("unknown directory")
	ErrUnknownVersion   = errors.New("unknwon version")
	ErrAmbiguousVersion = errors.New("ambiguous version")
)
