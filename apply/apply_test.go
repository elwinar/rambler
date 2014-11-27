package apply

import (
	"database/sql"
	"errors"
	. "github.com/franela/goblin"
	"testing"
)

func TestApply(t *testing.T) {
	g := Goblin(t)
	
	var migration MockMigration = MockMigration {
		statements: map[string][]string {
			"up": []string{
				"one",
				"three",
			},
			"down": []string{
				"two",
				"four",
			},
		},
	}

	var execs int
	var commits int
	var rollbacks int

	g.Describe("Apply", func() {
		g.BeforeEach(func() {
			execs = 0
			commits = 0
			rollbacks = 0
		})

		g.It("Should return an error on nil migration", func() {
			err, sqlerr := apply(nil, MockTxer{})
			g.Assert(err).Equal(ErrNilMigration)
			g.Assert(sqlerr).Equal(nil)
		})

		g.It("Should execute migration's up statements in order", func() {
			var index int = 0
			var fail bool = false
			tx := MockTxer{
				exec: func(query string, args ...interface{}) (sql.Result, error) {
					if query != migration.statements["up"][index] {
						fail = true
					}
					index++
					execs++
					return MockResult{}, nil
				},
				commit: func() error {
					commits++
					return nil
				},
				rollback: func() error {
					rollbacks++
					return nil
				},
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
			tx := MockTxer{
				exec: func(query string, args ...interface{}) (sql.Result, error) {
					execs++
					return MockResult{}, errors.New("error")
				},
				commit: func() error {
					commits++
					return nil
				},
				rollback: func() error {
					rollbacks++
					return nil
				},
			}

			err, sqlerr := apply(&migration, tx)
			g.Assert(execs).Equal(1)
			g.Assert(commits).Equal(0)
			g.Assert(rollbacks).Equal(1)
			g.Assert(err).Equal(nil)
			g.Assert(sqlerr).Equal(errors.New("error"))
		})
		
		g.It("Should return an error on commit fail", func(){
			tx := MockTxer{
				exec: func(query string, args ...interface{}) (sql.Result, error) {
					execs++
					return MockResult{}, nil
				},
				commit: func() error {
					commits++
					return errors.New("error")
				},
				rollback: func() error {
					rollbacks++
					return nil
				},
			}

			err, sqlerr := apply(&migration, tx)
			g.Assert(execs).Equal(2)
			g.Assert(commits).Equal(1)
			g.Assert(rollbacks).Equal(0)
			g.Assert(err).Equal(errors.New("error"))
			g.Assert(sqlerr).Equal(nil)
		})
		
		g.It("Should return an error on rollback fail", func(){
			tx := MockTxer{
				exec: func(query string, args ...interface{}) (sql.Result, error) {
					execs++
					return MockResult{}, errors.New("error")
				},
				commit: func() error {
					commits++
					return nil
				},
				rollback: func() error {
					rollbacks++
					return errors.New("error")
				},
			}

			err, sqlerr := apply(&migration, tx)
			g.Assert(execs).Equal(1)
			g.Assert(commits).Equal(0)
			g.Assert(rollbacks).Equal(1)
			g.Assert(err).Equal(errors.New("error"))
			g.Assert(sqlerr).Equal(errors.New("error"))
		})
		
		g.It("Should return nil on success", func(){
			tx := MockTxer{
				exec: func(query string, args ...interface{}) (sql.Result, error) {
					execs++
					return MockResult{}, nil
				},
				commit: func() error {
					commits++
					return nil
				},
				rollback: func() error {
					rollbacks++
					return nil
				},
			}

			err, sqlerr := apply(&migration, tx)
			g.Assert(execs).Equal(2)
			g.Assert(commits).Equal(1)
			g.Assert(rollbacks).Equal(0)
			g.Assert(err).Equal(nil)
			g.Assert(sqlerr).Equal(nil)
		})
	})
}
