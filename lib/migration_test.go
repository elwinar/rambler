package lib

import (
	"errors"
	. "github.com/franela/goblin"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

var (
	nilMigration *Migration
)

func TestNewMigration(t *testing.T) {
	g := Goblin(t)
	
	g.Describe("NewMigration", func() {
		g.It("Should reject unknown directory path", func() {
			m, err := newMigration("", 0, func(pattern string) ([]string, error) {
				return nil, errors.New("unknown")
			}, func(path string) (io.ReadSeeker, error) {
				return nil, nil
			})
			g.Assert(err).Equal(ErrUnknownDirectory)
			g.Assert(m).Equal(nilMigration)
		})

		g.It("Should reject unknown migrations", func() {
			m, err := newMigration("", 0, func(pattern string) ([]string, error) {
				return nil, nil
			}, func(path string) (io.ReadSeeker, error) {
				return nil, nil
			})
			g.Assert(err).Equal(ErrUnknownVersion)
			g.Assert(m).Equal(nilMigration)
		})

		g.It("Should reject ambiguous migrations", func() {
			m, err := newMigration("", 0, func(pattern string) ([]string, error) {
				return []string{"a", "b"}, nil
			}, func(path string) (io.ReadSeeker, error) {
				return nil, nil
			})
			g.Assert(err).Equal(ErrAmbiguousVersion)
			g.Assert(m).Equal(nilMigration)
		})

		g.It("Should parse filenames to get descriptions", func() {
			m, err := newMigration("", 0, func(pattern string) ([]string, error) {
				return []string{"42_forty_two.sql"}, nil
			}, func(path string) (io.ReadSeeker, error) {
				return nil, nil
			})
			g.Assert(err).Equal(nil)
			g.Assert(m.Description).Equal("forty_two")
		})

		g.It("Should open a handle for the migration file", func() {
			m, err := newMigration("", 0, func(pattern string) ([]string, error) {
				return []string{"42_forty_two.sql"}, nil
			}, func(path string) (io.ReadSeeker, error) {
				return strings.NewReader("rambler"), nil
			})
			g.Assert(err).Equal(nil)

			content, err := ioutil.ReadAll(m.reader)
			g.Assert(err).Equal(nil)
			g.Assert(content).Equal([]byte("rambler"))
		})
	})
}

func TestScan(t *testing.T) {
	g := Goblin(t)
	
	var text string = `
-- rambler up
one
-- rambler down
two
-- rambler up
three
-- rambler down
four
`
	r := strings.NewReader(text)
	
	var reads int
	var bytes int
	var seeks int
	
	m := &Migration{
		reader: &MockReader{
			seek: func(offset int64, whence int) (int64, error) {
				seeks++
				return r.Seek(offset, whence)
			},
			read: func(p []byte) (int, error) {
				reads++
				b, err := r.Read(p)
				bytes += b
				return bytes, err
			},
		},
	}
	
	g.Describe("Scan", func() {
		g.BeforeEach(func() {
			reads = 0
			bytes = 0
			seeks = 0
		})
		
		g.It("Should rewind the reader", func() {
			m.Scan("up")
			g.Assert(seeks).Equal(1)
		})
		
		g.It("Should read the whole file", func() {
			m.Scan("up")
			g.Assert(bytes).Equal(len(text))
		})

		g.It("Should find statements if there is statements to find", func() {
			statements := m.Scan("up")
			g.Assert(statements).Equal([]string{"one", "three"})
		})

		g.It("Should return an empty slice if there is no statements to find", func() {
			statements := m.Scan("right")
			g.Assert(len(statements)).Equal(0)
		})
	})
}
