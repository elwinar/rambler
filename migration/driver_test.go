package migration

import (
	. "github.com/franela/goblin"
	"testing"
)

type MockDriver struct{}

func TestRegister(t *testing.T) {
	g := Goblin(t)
	g.Describe("Register", func() {
		g.It("Should accept drivers", func() {
			drivers := make(map[string]Driver)
			err := register("mock", MockDriver{}, drivers)
			g.Assert(err).Equal(nil)
		})
		
		g.It("Should not accept the same driver twice", func() {
			drivers := make(map[string]Driver)
			err := register("mock", MockDriver{}, drivers)
			g.Assert(err).Equal(nil)
			
			err = register("mock", MockDriver{}, drivers)
			g.Assert(err).Equal(ErrDriverAlreadyRegistered)
		})
	})
}
