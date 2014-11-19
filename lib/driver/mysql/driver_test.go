package mysql

import (
	. "github.com/franela/goblin"
	"github.com/jmoiron/sqlx"
	"testing"
)

var (
	nilDriver *Driver
)

func TestNewDriver(t *testing.T) {
	g := Goblin(t)
	g.Describe("NewDriver", func() {
		g.It("Should reject invalid DSN", func() {
			d, err := newDriver("invalid", sqlx.Connect)
			g.Assert(err).Equal(ErrUnknownDatabase)
			g.Assert(d).Equal(nilDriver)
		})

		g.It("Should fail on unreachable database", func() {
			d, err := newDriver("/unreachable", sqlx.Connect)
			g.Assert(err).Equal(ErrUnknownDatabase)
			g.Assert(d).Equal(nilDriver)
		})

		g.It("Should parse the DSN to get the schema", func() {
			d, err := newDriver("root:@tcp(localhost/mysql:3306)/schema?parseTime=true", func(driver, dsn string) (*sqlx.DB, error) {
				return &sqlx.DB{}, nil
			})
			g.Assert(err).Equal(nil)
			g.Assert(d.schema).Equal("schema")
		})
	})
}
