package rambler

import (
	"errors"
	. "github.com/franela/goblin"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	unknownDirectory = "unknown-dir/"
	unknownVersion   = 13

	knownPath        = "test/42_forty_two.sql"
	knownDirectory   = "test/"
	knownVersion     = 42
	knownDescription = "forty_two"
	knownContent     = `-- rambler up
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

func TestNewMigration(t *testing.T) {
	g := Goblin(t)
	g.Describe("NewMigration", func() {
		g.It("Should reject unknown directory path", func() {
			m, err := newMigration(unknownDirectory, knownVersion, func(pattern string) ([]string, error) {
				return nil, errors.New("unknown")
			}, func(path string) (io.ReadSeeker, error) {
				return nil, nil
			})
			g.Assert(err).Equal(ErrUnknownDirectory)
			g.Assert(m).Equal(nilMigration)
		})

		g.It("Should reject unknown migrations", func() {
			m, err := newMigration(knownDirectory, unknownVersion, func(pattern string) ([]string, error) {
				return nil, nil
			}, func(path string) (io.ReadSeeker, error) {
				return nil, nil
			})
			g.Assert(err).Equal(ErrUnknownVersion)
			g.Assert(m).Equal(nilMigration)
		})

		g.It("Should reject ambiguous migrations", func() {
			m, err := newMigration(knownDirectory, ambiguousVersion, func(pattern string) ([]string, error) {
				return []string{"a", "b"}, nil
			}, func(path string) (io.ReadSeeker, error) {
				return nil, nil
			})
			g.Assert(err).Equal(ErrAmbiguousVersion)
			g.Assert(m).Equal(nilMigration)
		})

		g.It("Should parse filenames to get descriptions", func() {
			m, err := newMigration(knownDirectory, knownVersion, func(pattern string) ([]string, error) {
				return []string{knownPath}, nil
			}, func(path string) (io.ReadSeeker, error) {
				return nil, nil
			})
			g.Assert(err).Equal(nil)
			g.Assert(m.Description).Equal(knownDescription)
		})

		g.It("Should open a handle for the migration file", func() {
			m, err := newMigration(knownDirectory, knownVersion, func(pattern string) ([]string, error) {
				return []string{knownPath}, nil
			}, func(path string) (io.ReadSeeker, error) {
				return strings.NewReader(knownContent), nil
			})
			g.Assert(err).Equal(nil)

			content, err := ioutil.ReadAll(m.reader)
			g.Assert(err).Equal(nil)
			g.Assert(content).Equal([]byte(knownContent))
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
				Reader:  strings.NewReader(knownContent),
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
				Reader:  strings.NewReader(knownContent),
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
