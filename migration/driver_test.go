package migration

import (
	"github.com/elwinar/rambler/configuration"
	. "github.com/franela/goblin"
	"testing"
)

func TestRegisterDriver(t *testing.T) {
	g := Goblin(t)

	var constructor func(configuration.Environment) (Driver, error)
	var constructors map[string]Constructor

	g.Describe("Register", func() {
		g.BeforeEach(func() {
			constructor = func(configuration.Environment) (Driver, error) {
				return nil, nil
			}

			constructors = make(map[string]Constructor)
		})

		g.It("Should register new drivers", func() {
			err := registerDriver("mock", constructor, constructors)
			g.Assert(err).Equal(nil)
			g.Assert(len(constructors)).Equal(1)
		})

		g.It("Shouldn't accept the same driver twice", func() {
			err := registerDriver("mock", constructor, constructors)
			g.Assert(err).Equal(nil)

			err = registerDriver("mock", constructor, constructors)
			g.Assert(err).Equal(ErrDriverAlreadyRegistered)
		})
	})
}

func TestGetDriver(t *testing.T) {
	g := Goblin(t)

	var constructor func(configuration.Environment) (Driver, error)
	var constructors map[string]Constructor
	var env configuration.Environment
	var driver MockDriver

	g.Describe("Get", func() {
		g.BeforeEach(func() {
			constructor = func(configuration.Environment) (Driver, error) {
				return &driver, nil
			}

			constructors = make(map[string]Constructor)

			env.Driver = "mock"
		})

		g.It("Should retrieve existing drivers", func() {
			constructors["mock"] = constructor

			d, err := getDriver(env, constructors)
			g.Assert(err).Equal(nil)
			g.Assert(d).Equal(&driver)
		})

		g.It("Should fail on unknown driver", func() {
			d, err := getDriver(env, constructors)
			g.Assert(d).Equal(nil)
			g.Assert(err).Equal(ErrDriverNotRegistered)
		})
	})
}
