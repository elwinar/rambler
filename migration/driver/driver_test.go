package driver

import (
	. "github.com/franela/goblin"
	"testing"
)

type MockDriver struct{}

func (d MockDriver) MigrationTableExists() (bool, error) {
	return false, nil
}

func (d MockDriver) CreateMigrationTable() (error) {
	return nil
}

func MockConstructor(options string) (Driver, error) {
	return MockDriver{}, nil
}

func TestRegisterDriver(t *testing.T) {
	g := Goblin(t)
	g.Describe("Register", func() {
		g.It("Should register new drivers", func() {
			constructors := make(map[string]Constructor)
			
			err := register("mock", MockConstructor, constructors)
			g.Assert(err).Equal(nil)
			g.Assert(len(constructors)).Equal(1)
		})
		
		g.It("Shouldn't accept the same driver twice", func() {
			constructors := make(map[string]Constructor)
			
			err := register("mock", MockConstructor, constructors)
			g.Assert(err).Equal(nil)
			
			err = register("mock", MockConstructor, constructors)
			g.Assert(err).Equal(ErrDriverAlreadyRegistered)
		})
	})
}

func TestGetDriver(t *testing.T) {
	g := Goblin(t)
	g.Describe("Get", func() {
		g.It("Should retrieve existing drivers", func() {
			constructors := make(map[string]Constructor)
			constructors["mock"] = MockConstructor
			
			driver, err := get("mock", "", constructors)
			g.Assert(err).Equal(nil)
			g.Assert(driver).Equal(MockDriver{})
		})
		
		g.It("Should fail on unknown driver", func() {
			drivers := make(map[string]Constructor)
			
			driver, err := get("mock", "", drivers)
			g.Assert(driver).Equal(nil)
			g.Assert(err).Equal(ErrDriverNotRegistered)
		})
	})
}
