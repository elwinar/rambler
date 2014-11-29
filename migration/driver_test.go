package migration

import (
	"github.com/elwinar/rambler/configuration"
	. "github.com/franela/goblin"
	"testing"
)

type MockDriver struct{}

func (d MockDriver) MigrationTableExists() (bool, error) {
	return false, nil
}

func (d MockDriver) CreateMigrationTable() error {
	return nil
}

func MockConstructor(env configuration.Environment) (Driver, error) {
	return MockDriver{}, nil
}

func TestRegisterDriver(t *testing.T) {
	g := Goblin(t)
	g.Describe("Register", func() {
		g.It("Should registerDriver new drivers", func() {
			constructors := make(map[string]Constructor)

			err := registerDriver("mock", MockConstructor, constructors)
			g.Assert(err).Equal(nil)
			g.Assert(len(constructors)).Equal(1)
		})

		g.It("Shouldn't accept the same driver twice", func() {
			constructors := make(map[string]Constructor)

			err := registerDriver("mock", MockConstructor, constructors)
			g.Assert(err).Equal(nil)

			err = registerDriver("mock", MockConstructor, constructors)
			g.Assert(err).Equal(ErrDriverAlreadyRegistered)
		})
	})
}

func TestGetDriver(t *testing.T) {
	g := Goblin(t)
	
	mockEnv := configuration.Environment{
		Driver: "mock",
	}
	
	g.Describe("Get", func() {
		g.It("Should retrieve existing drivers", func() {
			constructors := make(map[string]Constructor)
			constructors["mock"] = MockConstructor

			driver, err := getDriver(mockEnv, constructors)
			g.Assert(err).Equal(nil)
			g.Assert(driver).Equal(MockDriver{})
		})

		g.It("Should fail on unknown driver", func() {
			drivers := make(map[string]Constructor)

			driver, err := getDriver(mockEnv, drivers)
			g.Assert(driver).Equal(nil)
			g.Assert(err).Equal(ErrDriverNotRegistered)
		})
	})
}
