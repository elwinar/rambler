package migration

import (
	"errors"
	"github.com/elwinar/rambler/configuration"
	. "github.com/franela/goblin"
	"os"
	"testing"
)

func TestNewService(t *testing.T) {
	g := Goblin(t)

	var nilservice *service
	var env configuration.Environment = configuration.Environment{
		Driver:    "mock",
		Directory: "dir",
	}

	var stat func(dir string) (os.FileInfo, error)
	var driverConstructor func(env configuration.Environment) (Driver, error)

	g.Describe("NewService", func() {
		g.BeforeEach(func() {
			stat = func(dir string) (os.FileInfo, error) {
				return nil, nil
			}

			driverConstructor = func(env configuration.Environment) (Driver, error) {
				return nil, nil
			}
		})

		g.It("Should reject unknown directory path", func() {
			s, err := newService(env, func(dir string) (os.FileInfo, error) {
				return nil, errors.New("error")
			}, driverConstructor)
			g.Assert(err).Equal(ErrUnknownDirectory)
			g.Assert(s).Equal(nilservice)
		})

		g.It("Should reject unknown driver", func() {
			s, err := newService(env, stat, func(env configuration.Environment) (Driver, error) {
				return nil, errors.New("error")
			})
			g.Assert(err).Equal(ErrUnknownDriver)
			g.Assert(s).Equal(nilservice)
		})

		g.It("Should return an initialized service", func() {
			s, err := newService(env, stat, driverConstructor)
			g.Assert(err).Equal(nil)
			g.Assert(s.env).Equal(env)
			g.Assert(s.driver).Equal(nil)
		})
	})
}
