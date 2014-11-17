package rambler

import (
	. "github.com/franela/goblin"
	"io/ioutil"
	"testing"
)

const (
	unknownDirectory = "unknown-dir/"
	unknownVersion = 13
	
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
			g.Assert(m.Description).Equal(knownDescription)
			g.Assert(err).Equal(nil)
		})
		
		g.It("Should get the right migration based on version number", func() {
			m, err := NewMigration(knownDirectory, knownVersion)
			content, err := ioutil.ReadAll(m.reader)
			g.Assert(content).Equal([]byte(knownContent))
			g.Assert(err).Equal(nil)
		})
	})
}
