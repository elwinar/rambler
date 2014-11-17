package rambler

import (
	"io"
)

type Opener interface {
	Open(string) (io.ReadSeeker, error)
}
