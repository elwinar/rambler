package lib

import (
	"errors"
)

// The various errors returned by the package
var (
	ErrUnknownDirectory   = errors.New("unknown directory")
	ErrUnknownVersion     = errors.New("unknwon version")
	ErrUnknownDriver      = errors.New("unknwon driver")
	ErrUnknownEnvironment = errors.New("unknwon environment")
	ErrAmbiguousVersion   = errors.New("ambiguous version")
)
