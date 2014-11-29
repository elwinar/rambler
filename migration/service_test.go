package migration

import (
	"errors"
	"github.com/elwinar/rambler/driver"
	. "github.com/franela/goblin"
	"os"
	"testing"
)

var (
	nilService *service
)

type MockDriver struct{}

func (d MockDriver) MigrationTableExists() (bool, error) {
	return false, nil
}

func (d MockDriver) CreateMigrationTable() error {
	return nil
}

func TestNewService(t *testing.T) {
	g := Goblin(t)
	g.Describe("NewService", func() {
		g.It("Should reject unknown directory path", func() {
			s, err := newService("", "", "", func(dir string) (os.FileInfo, error) {
				return nil, errors.New("error")
			}, func(string, string) (driver.Driver, error) {
				return nil, nil
			})
			g.Assert(err).Equal(ErrUnknownDirectory)
			g.Assert(s).Equal(nilService)
		})

		g.It("Should reject unknown driver", func() {
			s, err := newService("", "", "", func(dir string) (os.FileInfo, error) {
				return nil, nil
			}, func(string, string) (driver.Driver, error) {
				return nil, errors.New("error")
			})
			g.Assert(err).Equal(ErrUnknownDriver)
			g.Assert(s).Equal(nilService)
		})

		g.It("Should return an initialized service", func() {
			d := &MockDriver{}
			s, err := newService("mock", "", "dir", func(dir string) (os.FileInfo, error) {
				return nil, nil
			}, func(string, string) (driver.Driver, error) {
				return d, nil
			})
			g.Assert(err).Equal(nil)
			g.Assert(s.directory).Equal("dir")
			g.Assert(s.driver).Equal(d)
		})
	})
}
