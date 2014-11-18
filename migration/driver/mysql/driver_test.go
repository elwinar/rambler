package mysql

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestConstructor(t *testing.T) {
	g := Goblin(t)
	g.Describe("Constructor", func() {
		g.It("Should reject invalid go-sql-driver DSN")
		g.It("Should fail on unreachable database")
		g.It("Should initialize the database handle")
	})
}
