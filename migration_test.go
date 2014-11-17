package rambler

import (
	. "github.com/franela/goblin"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	unknownDirectory = "unknown-dir/"
	unknownVersion = 13
	
	knownPath = "test/42_forty_two.sql"
	knownDirectory = "test/"
	knownVersion = 42
	knownDescription = "forty_two"
	knownContent = `-- rambler up
CREATE TABLE foo (
	id INTEGER UNSIGNED AUTO_INCREMENT,
	PRIMARY KEY (id)
);

-- rambler down
DROP TABLE foo;
`
	
	ambiguousVersion = 33
)

var (
	nilMigration *Migration
)

type MockGlober struct {
	Glober
	cb      func(string) ([]string, error)
}

func (g MockGlober) Glob(pattern string) ([]string, error) {
	return g.cb(pattern)
}

type MockOpener struct {
	Opener
	cb      func(string) (io.ReadSeeker, error)
}

func (o MockOpener) Open(path string) (io.ReadSeeker, error) {
	return o.cb(path)
}

func TestNewMigration(t *testing.T) {
	g := Goblin(t)
	g.Describe("NewMigration", func() {
		g.It("Should reject unknown directory path", func() {
			m, err := NewMigration(unknownDirectory, knownVersion, MockGlober{
				cb: func(pattern string) ([]string, error) {
					return nil, errors.New("unknown")
				},
			}, MockOpener{
				cb: func(path string) (io.ReadSeeker, error) {
					return nil, nil
				},
			})
			g.Assert(m).Equal(nilMigration)
			g.Assert(err).Equal(ErrUnknownDirectory)
		})
		
		g.It("Should reject unknown migrations", func() {
			m, err := NewMigration(knownDirectory, unknownVersion, MockGlober{
				cb: func(pattern string) ([]string, error) {
					return nil, nil
				},
			}, MockOpener{
				cb: func(path string) (io.ReadSeeker, error) {
					return nil, nil
				},
			})
			g.Assert(m).Equal(nilMigration)
			g.Assert(err).Equal(ErrUnknownVersion)
		})
		
		g.It("Should reject ambiguous migrations", func() {
			m, err := NewMigration(knownDirectory, ambiguousVersion, MockGlober{
				cb: func(pattern string) ([]string, error) {
					return []string{"a","b"}, nil
				},
			}, MockOpener{
				cb: func(path string) (io.ReadSeeker, error) {
					return nil, nil
				},
			})
			g.Assert(m).Equal(nilMigration)
			g.Assert(err).Equal(ErrAmbiguousVersion)
		})
		
		g.It("Should parse filenames to get descriptions", func() {
			m, err := NewMigration(knownDirectory, knownVersion, MockGlober{
				cb: func(pattern string) ([]string, error) {
					return []string{knownPath}, nil
				},
			}, MockOpener{
				cb: func(path string) (io.ReadSeeker, error) {
					return nil, nil
				},
			})
			g.Assert(m.Description).Equal(knownDescription)
			g.Assert(err).Equal(nil)
		})
		
		g.It("Should open a handle for the migration file", func() {
			m, err := NewMigration(knownDirectory, knownVersion, MockGlober{
				cb: func(pattern string) ([]string, error) {
					return []string{knownPath}, nil
				},
			}, MockOpener{
				cb: func(path string) (io.ReadSeeker, error) {
					return strings.NewReader(knownContent), nil
				},
			})
			g.Assert(err).Equal(nil)
			
			content, err := ioutil.ReadAll(m.reader)
			g.Assert(content).Equal([]byte(knownContent))
			g.Assert(err).Equal(nil)
		})
	})
}

type MockSeeker struct {
	*strings.Reader
	counter int
}

func (s *MockSeeker) Seek(offset int64, whence int) (int64, error) {
	s.counter++
	return s.Reader.Seek(offset, whence)
}

type MockReader struct {
	*strings.Reader
	counter int
}

func (r *MockReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.counter += n
	return n, err
}

func TestScan(t *testing.T) {
	g := Goblin(t)
	g.Describe("Scan", func() {
		g.It("Should rewind the reader", func() {
			seeker := &MockSeeker{
				Reader: strings.NewReader(knownContent),
				counter: 0,
			}
			m := &Migration{
				reader: seeker,
			}
			m.Scan("up")
			g.Assert(seeker.counter).Equal(1)
		})
		
		g.It("Should read the whole file", func() {
			reader := &MockReader{
				Reader: strings.NewReader(knownContent),
				counter: 0,
			}
			m := &Migration{
				reader: reader,
			}
			m.Scan("up")
			g.Assert(reader.counter).Equal(len([]byte(knownContent)))
		})
	})
}
