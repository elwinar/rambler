package apply

/*
import (
	"errors"
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration"
	. "github.com/franela/goblin"
	"testing"
)

func TestCommand(t *testing.T) {
	g := Goblin(t)
	
	var env = configuration.Environment{
		Driver:    "mock",
		Directory: "dir",
	}
	
	var tx MockTransaction

	var s MockService
	var exists int
	var creates int
	var listApplieds int
	var listAvailables int
	var starts int
	
	var newService serviceConstructor
	var news int
	
	var f filter
	var filters int
	
	var m MockMigration
	var newMigration migrationConstructor
	var newMigrations int
	
	var a applier
	var applies int

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
			
			s.startTransaction = func() (migration.Transaction, error) {
				starts++
				return tx, nil
			}

			exists = 0
			creates = 0
			listApplieds = 0
			listAvailables = 0
			starts = 0
			
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
			
			m.scan = func(_ string) []string {
				return nil
			}
			
			newMigration = func(_ string, _ uint64) (scanner, error) {
				newMigrations++
				return &m, nil
			}
			
			newMigrations = 0
			
			a = func(_ scanner, _ txer) (error, error) {
				applies++
				return nil, nil
			}
			
			applies = 0
		})

		g.It("Should fail on invalid environment", func() {
			newService = func(env configuration.Environment) (migration.Service, error) {
				return nil, errors.New("new service error")
			}
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("new service error"))
		})
		
		g.It("Should check for the migration table", func() {
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(nil)
			g.Assert(exists).Equal(1)
		})
		
		g.It("Should return an error if an error occurs while checking for the migration table", func() {
			s.migrationTableExists = func() (bool, error) {
				exists++
				return false, errors.New("exists error")
			}
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("exists error"))
			g.Assert(exists).Equal(1)
			g.Assert(creates).Equal(0)
		})
		
		g.It("Should create the migration table if it does'nt exists", func() {
			s.migrationTableExists = func() (bool, error) {
				exists++
				return false, nil
			}
			err := command(env, false, newService, f, newMigration, a)
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
				return errors.New("create error")
			}
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("create error"))
			g.Assert(exists).Equal(1)
			g.Assert(creates).Equal(1)
			g.Assert(listApplieds).Equal(0)
		})
		
		g.It("Shouldn't create the migration table if it already exists", func() {
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(nil)
			g.Assert(creates).Equal(0)
		})
		
		g.It("Should list the already applied migrations", func() {
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(nil)
			g.Assert(listApplieds).Equal(1)
		})
		
		g.It("Should return an error if an error occurs while listing already applied migrations", func() {
			s.listAppliedMigrations = func() ([]uint64, error) {
				listApplieds++
				return nil, errors.New("list applied error")
			}
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("list applied error"))
			g.Assert(listApplieds).Equal(1)
		})
		
		g.It("Should list the available migrations", func() {
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(nil)
			g.Assert(listAvailables).Equal(1)
		})
		
		g.It("Should return an error if an error occurs while listing available migrations", func() {
			s.listAvailableMigrations = func() ([]uint64, error) {
				listAvailables++
				return nil, errors.New("list available error")
			}
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("list available error"))
			g.Assert(listAvailables).Equal(1)
		})
		
		g.It("Should filter out the migrations already applied", func() {
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(nil)
			g.Assert(filters).Equal(1)
		})
		
		g.It("Should return an error if an error occurs while filtering migrations", func() {
			f = func([]uint64, []uint64) ([]uint64, error) {
				filters++
				return nil, errors.New("filter error")
			}
			
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("filter error"))
			g.Assert(filters).Equal(1)
		})
		
		g.It("Should return an error if an error occurs while retrieving a migration", func() {
			newMigration = func(_ string, _ uint64) (scanner, error) {
				return &m, errors.New("new migration error")
			}
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("new migration error"))
		})
		
		g.It("Should return an error if an error occurs while starting a transaction", func() {
			s.startTransaction = func() (migration.Transaction, error) {
				return nil, errors.New("start transaction error")
			}
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("start transaction error"))
		})
		
		g.It("Should apply one migration if requested", func() {
			f = func([]uint64, []uint64) ([]uint64, error) {
				filters++
				return []uint64{1, 2, 3}, nil
			}
			
			err := command(env, false, newService, f, newMigration, a)
			g.Assert(err).Equal(nil)
			g.Assert(applies).Equal(1)
		})
		
		g.It("Should apply all migration if requested", func() {
			f = func([]uint64, []uint64) ([]uint64, error) {
				filters++
				return []uint64{1, 2, 3}, nil
			}
			
			err := command(env, true, newService, f, newMigration, a)
			g.Assert(err).Equal(nil)
			g.Assert(applies).Equal(3)
		})
		
		g.It("Should return an error if an error occurs while appling a migration", func() {
			a = func(_ scanner, _ txer) (error, error) {
				applies++
				return errors.New("apply sql error"), nil
			}
			
			err := command(env, true, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("apply sql error"))
			
			a = func(_ scanner, _ txer) (error, error) {
				applies++
				return errors.New("apply sql error"), errors.New("apply tx error")
			}
			
			err = command(env, true, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("apply sql error; apply tx error"))
			
			a = func(_ scanner, _ txer) (error, error) {
				applies++
				return nil, errors.New("apply tx error")
			}
			
			err = command(env, true, newService, f, newMigration, a)
			g.Assert(err).Equal(errors.New("apply tx error"))
		})
	})
}
*/
