package apply 

import (
	. "github.com/franela/goblin"
	"testing"
)

type MockService struct {}

func TestCommand(t *testing.T) {
	g := Goblin(t)
	
	g.Describe("Command", func() {
		g.It("Should check for the migration table")
		g.It("Should create the migration table if it does'nt exists")
		g.It("Should list the already applied migrations")
		g.It("Should filter out the migrations already applied")
		g.It("Should apply one migration if requested")
		g.It("Should apply all migrations in order if requested")
		g.It("Should stop on error")
	})
}
