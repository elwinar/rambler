package apply

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestCommand(t *testing.T) {
	g := Goblin(t)

	var s MockService
	var exists int
	var creates int

	g.Describe("Command", func() {

		g.BeforeEach(func() {
			s.migrationTableExists = func() (bool, error) {
				exists++
				return true, nil
			}

			s.createMigrationTable = func() error {
				creates++
				return nil
			}

			exists = 0
			creates = 0
		})

		g.It("Should initialize a new service", func() {

		})
		g.It("Should check for the migration table")
		g.It("Should create the migration table if it does'nt exists")
		g.It("Should list the already applied migrations")
		g.It("Should filter out the migrations already applied")
		g.It("Should apply one migration if requested")
		g.It("Should apply all migrations in order if requested")
		g.It("Should stop on error")
	})
}
