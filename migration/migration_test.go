package migration

import (
	"errors"
	. "github.com/franela/goblin"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewMigration(t *testing.T) {
	g := Goblin(t)

	var glob glober
	var open opener

	var nilmigration *Migration

	g.Describe("NewMigration", func() {
		g.BeforeEach(func() {
			glob = func(string) ([]string, error) {
				return []string{"1_first.sql"}, nil
			}

			open = func(string) (io.ReadSeeker, error) {
				return nil, nil
			}
		})

		g.It("reject unknown directory path", func() {
			glob = func(pattern string) ([]string, error) {
				return nil, errors.New("unknown")
			}
			
			m, err := newMigration("", 0, glob, open)
			g.Assert(err).Equal(ErrUnknownDirectory)
			g.Assert(m).Equal(nilmigration)
		})

		g.It("reject unknown migrations", func() {
			glob = func(string) ([]string, error) {
				return nil, nil
			}

			m, err := newMigration("", 0, glob, open)
			g.Assert(err).Equal(ErrUnknownVersion)
			g.Assert(m).Equal(nilmigration)
		})

		g.It("reject ambiguous migrations", func() {
			glob = func(pattern string) ([]string, error) {
				return []string{"a", "b"}, nil
			}
			m, err := newMigration("", 0, glob, open)
			g.Assert(err).Equal(ErrAmbiguousVersion)
			g.Assert(m).Equal(nilmigration)
		})

		g.It("open a handle for the migration file", func() {
			open = func(path string) (io.ReadSeeker, error) {
				return strings.NewReader("rambler"), nil
			}
			m, err := newMigration("", 0, glob, open)
			g.Assert(err).Equal(nil)

			content, err := ioutil.ReadAll(m.reader)
			g.Assert(err).Equal(nil)
			g.Assert(content).Equal([]byte("rambler"))
		})
		
		g.It("fail on opening error", func() {
			open = func(_ string) (io.ReadSeeker, error) {
				return nil, errors.New("open error")
			}
			_, err := newMigration("", 0, glob, open)
			g.Assert(err).Equal(errors.New("open error"))
		})

		g.It("parse filenames to get descriptions", func() {
			m, err := newMigration("", 0, glob, open)
			g.Assert(err).Equal(nil)
			g.Assert(m.Description).Equal("first")
		})
	})
}

func TestScan(t *testing.T) {
	g := Goblin(t)

	var text = `
-- rambler up
one
-- rambler down
two
-- rambler up
three
-- rambler down
four
`
	var reader MockReader
	var migration = &Migration{
		reader: &reader,
	}
	var reads int
	var bytes int
	var seeks int

	g.Describe("Scan", func() {
		g.BeforeEach(func() {
			r := strings.NewReader(text)

			reader.seek = func(offset int64, whence int) (int64, error) {
				seeks++
				return r.Seek(offset, whence)
			}

			reader.read = func(p []byte) (int, error) {
				reads++
				b, err := r.Read(p)
				bytes += b
				return bytes, err
			}

			reads = 0
			bytes = 0
			seeks = 0
		})

		g.It("Should rewind the reader", func() {
			migration.Scan("up")
			g.Assert(seeks).Equal(1)
		})

		g.It("Should read the whole file", func() {
			migration.Scan("up")
			g.Assert(bytes).Equal(len(text))
		})

		g.It("Should find statements if there is statements to find", func() {
			statements := migration.Scan("up")
			g.Assert(statements).Equal([]string{"one", "three"})
		})

		g.It("Should return an empty slice if there is no statements to find", func() {
			statements := migration.Scan("right")
			g.Assert(len(statements)).Equal(0)
		})
		
		g.It("Read the last statement until the end", func() {
			statements := migration.Scan("down")
			g.Assert(len(statements)).Equal(2)
		})
	})
}
