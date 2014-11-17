package migration

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestNewService(t *testing.T) {
	g := Goblin(t)
	g.Describe("NewService", func() {
		g.It("Should use the correct service driver")
	})
}
