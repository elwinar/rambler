package rambler

import (
	. "github.com/franela/goblin"
	"strings"
	"testing"
)

const (
	unknownDirectory = "unknown-dir/"
	unknownVersion = 13
	
	knownDirectory = "test/"
	knownVersion = 42
	knownDescription = "forty_two"
	
	ambiguousVersion = 33
)

var (
	nilMigration *Migration
	knownMigration *Migration
)

func init() {
	knownMigration = &Migration{
		Version: knownVersion,
		Description: knownDescription,
		reader: strings.NewReader(""),
	}
}

func TestNewMigration(t *testing.T) {
	g := Goblin(t)
	g.Describe("NewMigration", func() {
		g.It("Should reject unknown directory path", func() {
			m, err := NewMigration(unknownDirectory, knownVersion)
			g.Assert(m).Equal(nilMigration)
			g.Assert(err).Equal(ErrUnknownDirectory)
		})
		
		g.It("Should reject unknown migrations", func() {
			m, err := NewMigration(knownDirectory, unknownVersion)
			g.Assert(m).Equal(nilMigration)
			g.Assert(err).Equal(ErrUnknownVersion)
		})
		
		g.It("Should reject ambiguous migrations", func() {
			m, err := NewMigration(knownDirectory, ambiguousVersion)
			g.Assert(m).Equal(nilMigration)
			g.Assert(err).Equal(ErrAmbiguousVersion)
		})
		
		g.It("Should parse filenames to get descriptions", func() {
			m, err := NewMigration(knownDirectory, knownVersion)
			g.Assert(m.Description).Equal(knownMigration.Description)
			g.Assert(err).Equal(nil)
		})
	})
}
