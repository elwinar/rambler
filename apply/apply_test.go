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
			sqlErr, txErr := apply(nil, MockTxer{})
			g.Assert(sqlErr).Equal(ErrNilMigration)
			g.Assert(txErr).Equal(nil)
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

			sqlErr, txErr := apply(&migration, tx)
			g.Assert(fail).Equal(false)
			g.Assert(execs).Equal(2)
			g.Assert(commits).Equal(1)
			g.Assert(rollbacks).Equal(0)
			g.Assert(sqlErr).Equal(nil)
			g.Assert(txErr).Equal(nil)
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

			sqlErr, txErr := apply(&migration, tx)
			g.Assert(execs).Equal(1)
			g.Assert(commits).Equal(0)
			g.Assert(rollbacks).Equal(1)
			g.Assert(sqlErr).Equal(errors.New("error"))
			g.Assert(txErr).Equal(nil)
		})
		
		g.It("Should return a transaction error on commit fail", func(){
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

			sqlErr, txErr := apply(&migration, tx)
			g.Assert(execs).Equal(2)
			g.Assert(commits).Equal(1)
			g.Assert(rollbacks).Equal(0)
			g.Assert(sqlErr).Equal(nil)
			g.Assert(txErr).Equal(errors.New("error"))
		})
		
		g.It("Should return a transaction error on rollback fail", func(){
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

			sqlErr, txErr := apply(&migration, tx)
			g.Assert(execs).Equal(1)
			g.Assert(commits).Equal(0)
			g.Assert(rollbacks).Equal(1)
			g.Assert(sqlErr).Equal(errors.New("error"))
			g.Assert(txErr).Equal(errors.New("error"))
		})
		
		g.It("Should return nil on success")
	})
}
