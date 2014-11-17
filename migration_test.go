package rambler

import (
	. "github.com/franela/goblin"
	"testing"
)

const (
	unknownDirectory = "unknown-dir/"
	knownDirectory = "test/"
	
	unknownVersion = 13
	knownVersion = 42
)

var (
	nilMigration *Migration
)

func TestNewMigration(t *testing.T) {
	g := Goblin(t)
	g.Describe("NewMigration", func() {
		g.It("Should reject unknown directory path", func() {
			m, err := NewMigration(unknownDirectory, knownVersion)
			g.Assert(m).Equal(nilMigration)
			g.Assert(err).Equal(ErrUnknownDirectory)
		})
	})
}
