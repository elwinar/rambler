package migration

import (
	"errors"
)

// The various errors returned by the package
var (
	ErrUnknownDirectory = errors.New("unknown directory")
	ErrUnknownVersion   = errors.New("unknwon version")
	ErrUnknownDriver    = errors.New("unknwon driver")
	ErrAmbiguousVersion = errors.New("ambiguous version")
)
