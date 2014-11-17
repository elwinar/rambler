package migration

import (
	. "github.com/franela/goblin"
	"testing"
)

type MockDriver struct{}

func TestRegisterDriver(t *testing.T) {
	g := Goblin(t)
	g.Describe("RegisterDriver", func() {
		g.It("Should register new drivers", func() {
			drivers := make(map[string]Driver)
			
			err := registerDriver("mock", MockDriver{}, drivers)
			g.Assert(err).Equal(nil)
			g.Assert(drivers["mock"]).Equal(MockDriver{})
		})
		
		g.It("Shouldn't accept the same driver twice", func() {
			drivers := make(map[string]Driver)
			
			err := registerDriver("mock", MockDriver{}, drivers)
			g.Assert(err).Equal(nil)
			
			err = registerDriver("mock", MockDriver{}, drivers)
			g.Assert(err).Equal(ErrDriverAlreadyRegistered)
		})
	})
}

func TestGetDriver(t *testing.T) {
	g := Goblin(t)
	g.Describe("GetDriver", func() {
		g.It("Should retrieve existing drivers", func() {
			drivers := make(map[string]Driver)
			drivers["mock"] = MockDriver{}
			
			driver, err := getDriver("mock", drivers)
			g.Assert(err).Equal(nil)
			g.Assert(driver).Equal(MockDriver{})
		})
		
		g.It("Should fail on unknown driver", func() {
			drivers := make(map[string]Driver)
			
			driver, err := getDriver("mock", drivers)
			g.Assert(driver).Equal(nil)
			g.Assert(err).Equal(ErrDriverNotRegistered)
		})
	})
}
