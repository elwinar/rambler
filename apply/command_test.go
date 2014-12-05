package apply

import (
	"errors"
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration"
	. "github.com/franela/goblin"
	"testing"
)

func TestCommand(t *testing.T) {
	g := Goblin(t)

	var s MockService
	var exists int
	var creates int
	
	var newService serviceConstructor
	var news int
	
	var env configuration.Environment = configuration.Environment{
		Driver:    "mock",
		Directory: "dir",
	}

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
			
			newService = func(env configuration.Environment) (migration.Service, error) {
				news++
				return s, nil
			}
			
			news = 0
		})

		g.It("Should fail on invalid environment", func() {
			err := command(env, false, func(env configuration.Environment) (migration.Service, error) {
				return s, errors.New(`error`)
			})
			g.Assert(err).Equal(errors.New(`error`))
		})
		
		g.It("Should check for the migration table", func() {
			err := command(env, false, newService)
			g.Assert(err).Equal(nil)
			g.Assert(exists).Equal(1)
		})
		
		g.It("Should create the migration table if it does'nt exists", func() {
			s.migrationTableExists = func() (bool, error) {
				exists++
				return false, nil
			}
			err := command(env, false, newService)
			g.Assert(err).Equal(nil)
			g.Assert(exists).Equal(1)
			g.Assert(creates).Equal(1)
		})
		
		g.It("Shouldn't create the migration table if it already exists", func() {
			err := command(env, false, newService)
			g.Assert(err).Equal(nil)
			g.Assert(exists).Equal(1)
			g.Assert(creates).Equal(0)
		})
		
		g.It("Should list the already applied migrations")
		
		g.It("Should filter out the migrations already applied")
		g.It("Should apply one migration if requested")
		g.It("Should apply all migrations in order if requested")
		g.It("Should stop on error")
	})
}
