package apply

import (
	"database/sql"
	"errors"
	. "github.com/franela/goblin"
	"testing"
)


func TestApply(t *testing.T) {
	g := Goblin(t)

	var migration MockMigration
	var scans int
	
	var tx MockTransaction
	var execs int
	var commits int
	var rollbacks int

	g.Describe("Apply", func() {
		g.BeforeEach(func() {
			// Re-initialize the migration mock
			migration.scan = func(_ string) []string {
				scans++
				return nil
			}
			
			scans = 0
			
			// Re-initialize the transaction mock
			tx.exec = func(_ string, _ ...interface{}) (sql.Result, error) {
				execs++
				return MockResult{}, nil
			}
			
			tx.commit = func() error {
				commits++
				return nil
			}
			
			tx.rollback = func() error {
				rollbacks++
				return nil
			}
			
			execs = 0
			commits = 0
			rollbacks = 0
		})

		g.It("Should return an error on nil migration", func() {
			err, sqlerr := apply(nil, tx)
			g.Assert(err).Equal(ErrNilMigration)
			g.Assert(sqlerr).Equal(nil)
		})

		g.It("Should execute migration's up statements in order", func() {
			var statements []string = []string{
				"one",
				"two",
			}
			var index int = 0
			var fail bool = false
			
			migration.scan = func(_ string) []string {
				return statements
			}
			
			tx.exec = func(query string, _ ...interface{}) (sql.Result, error) {
				if query != statements[index] {
					fail = true
				}
				index++
				execs++
				return MockResult{}, nil
			}

			err, sqlerr := apply(&migration, tx)
			g.Assert(fail).Equal(false)
			g.Assert(execs).Equal(2)
			g.Assert(commits).Equal(1)
			g.Assert(rollbacks).Equal(0)
			g.Assert(err).Equal(nil)
			g.Assert(sqlerr).Equal(nil)
		})

		g.It("Should rollback on SQL error", func() {
			migration.scan = func(_ string) []string {
				return []string{ "faulty" }
			}
			
			tx.exec = func(_ string, _ ...interface{}) (sql.Result, error) {
				execs++
				return MockResult{}, errors.New("error")
			}

			err, sqlerr := apply(&migration, tx)
			g.Assert(execs).Equal(1)
			g.Assert(commits).Equal(0)
			g.Assert(rollbacks).Equal(1)
			g.Assert(err).Equal(nil)
			g.Assert(sqlerr).Equal(errors.New("error"))
		})

		g.It("Should return an error on commit fail", func() {
			tx.commit = func() error {
				commits++
				return errors.New("error")
			}

			err, sqlerr := apply(&migration, tx)
			g.Assert(execs).Equal(0)
			g.Assert(commits).Equal(1)
			g.Assert(rollbacks).Equal(0)
			g.Assert(err).Equal(errors.New("error"))
			g.Assert(sqlerr).Equal(nil)
		})

		g.It("Should return an error on rollback fail", func() {
			migration.scan = func(_ string) []string {
				return []string{ "faulty" }
			}
			
			tx.exec = func(query string, args ...interface{}) (sql.Result, error) {
				execs++
				return MockResult{}, errors.New("error")
			}
			
			tx.rollback = func() error {
				rollbacks++
				return errors.New("error")
			}

			err, sqlerr := apply(&migration, tx)
			g.Assert(execs).Equal(1)
			g.Assert(commits).Equal(0)
			g.Assert(rollbacks).Equal(1)
			g.Assert(err).Equal(errors.New("error"))
			g.Assert(sqlerr).Equal(errors.New("error"))
		})

		g.It("Should return nil on success", func() {
			err, sqlerr := apply(&migration, tx)
			g.Assert(execs).Equal(0)
			g.Assert(commits).Equal(1)
			g.Assert(rollbacks).Equal(0)
			g.Assert(err).Equal(nil)
			g.Assert(sqlerr).Equal(nil)
		})
	})
}
