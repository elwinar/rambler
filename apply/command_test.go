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
	
	var env configuration.Environment = configuration.Environment{
		Driver:    "mock",
		Directory: "dir",
	}

	var s MockService
	var exists int
	var creates int
	var listApplieds int
	var listAvailables int
	
	var newService serviceConstructor
	var news int
	
	var f filter
	var filters int

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
			
			s.listAppliedMigrations = func() ([]uint64, error) {
				listApplieds++
				return nil, nil
			}
			
			s.listAvailableMigrations = func() ([]uint64, error) {
				listAvailables++
				return nil, nil
			}

			exists = 0
			creates = 0
			listApplieds = 0
			listAvailables = 0
			
			newService = func(_ configuration.Environment) (migration.Service, error) {
				news++
				return s, nil
			}
			
			news = 0
			
			f = func([]uint64, []uint64) ([]uint64, error) {
				filters++
				return []uint64{1, 2}, nil
			}
			
			filters = 0
		})

		g.It("Should fail on invalid environment", func() {
			newService = func(env configuration.Environment) (migration.Service, error) {
				return nil, errors.New(`error`)
			}
			err := command(env, false, newService, f)
			g.Assert(err).Equal(errors.New(`error`))
		})
		
		g.It("Should check for the migration table", func() {
			err := command(env, false, newService, f)
			g.Assert(err).Equal(nil)
			g.Assert(exists).Equal(1)
		})
		
		g.It("Should return an error if an error occurs while checking for the migration table", func() {
			s.migrationTableExists = func() (bool, error) {
				exists++
				return false, errors.New("error")
			}
			err := command(env, false, newService, f)
			g.Assert(err).Equal(errors.New("error"))
			g.Assert(exists).Equal(1)
			g.Assert(creates).Equal(0)
		})
		
		g.It("Should create the migration table if it does'nt exists", func() {
			s.migrationTableExists = func() (bool, error) {
				exists++
				return false, nil
			}
			err := command(env, false, newService, f)
			g.Assert(err).Equal(nil)
			g.Assert(exists).Equal(1)
			g.Assert(creates).Equal(1)
		})
		
		g.It("Should return an error if an error occurs while creating the migration table", func() {
			s.migrationTableExists = func() (bool, error) {
				exists++
				return false, nil
			}
			s.createMigrationTable = func() (error) {
				creates++
				return errors.New("error")
			}
			err := command(env, false, newService, f)
			g.Assert(err).Equal(errors.New("error"))
			g.Assert(exists).Equal(1)
			g.Assert(creates).Equal(1)
			g.Assert(listApplieds).Equal(0)
		})
		
		g.It("Shouldn't create the migration table if it already exists", func() {
			err := command(env, false, newService, f)
			g.Assert(err).Equal(nil)
			g.Assert(creates).Equal(0)
		})
		
		g.It("Should list the already applied migrations", func() {
			err := command(env, false, newService, f)
			g.Assert(err).Equal(nil)
			g.Assert(listApplieds).Equal(1)
		})
		
		g.It("Should return an error if an error occurs while listing already applied migrations", func() {
			s.listAppliedMigrations = func() ([]uint64, error) {
				listApplieds++
				return nil, errors.New("error")
			}
			err := command(env, false, newService, f)
			g.Assert(err).Equal(errors.New("error"))
			g.Assert(listApplieds).Equal(1)
		})
		
		g.It("Should list the available migrations", func() {
			err := command(env, false, newService, f)
			g.Assert(err).Equal(nil)
			g.Assert(listAvailables).Equal(1)
		})
		
		g.It("Should return an error if an error occurs while listing available migrations", func() {
			s.listAvailableMigrations = func() ([]uint64, error) {
				listAvailables++
				return nil, errors.New("error")
			}
			err := command(env, false, newService, f)
			g.Assert(err).Equal(errors.New("error"))
			g.Assert(listAvailables).Equal(1)
		})
		
		g.It("Should filter out the migrations already applied", func() {
			err := command(env, false, newService, f)
			g.Assert(err).Equal(nil)
			g.Assert(filters).Equal(1)
		})
		g.It("Should apply one migration if requested")
		g.It("Should apply all migrations in order if requested")
		g.It("Should stop on error")
	})
}
