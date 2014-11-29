package apply

import (
	"github.com/elwinar/rambler/migration"
	. "github.com/franela/goblin"
	"testing"
)

func TestFilter(t *testing.T) {
	g := Goblin(t)
	g.Describe("Filter", func() {
		g.It("Should complain about missing migrations", func() {
			available := []*migration.Migration{
				{Version: 1},
				{Version: 2},
				{Version: 4},
				{Version: 5},
			}
			applied := []*migration.Migration{
				{Version: 1},
				{Version: 2},
				{Version: 3},
				{Version: 4},
			}

			filtered, err := Filter(available, applied)
			g.Assert(err.Error()).Equal("missing migration 3")
			g.Assert(len(filtered)).Equal(0)
		})

		g.It("Should complain about out-of-order migrations", func() {
			available := []*migration.Migration{
				{Version: 1},
				{Version: 2},
				{Version: 3},
				{Version: 4},
			}
			applied := []*migration.Migration{
				{Version: 1},
				{Version: 2},
				{Version: 4},
				{Version: 5},
			}

			filtered, err := Filter(available, applied)
			g.Assert(err.Error()).Equal("out of order migration 3")
			g.Assert(len(filtered)).Equal(0)
		})

		g.It("Should return an empty slice when there is nothing to apply", func() {
			available := []*migration.Migration{
				{Version: 1},
				{Version: 2},
				{Version: 3},
				{Version: 4},
			}
			applied := []*migration.Migration{
				{Version: 1},
				{Version: 2},
				{Version: 3},
				{Version: 4},
			}

			filtered, err := Filter(available, applied)
			g.Assert(err).Equal(nil)
			g.Assert(len(filtered)).Equal(0)
		})

		g.It("Should return all remaining migrations where there is migrations to apply", func() {
			available := []*migration.Migration{
				{Version: 1},
				{Version: 2},
				{Version: 3},
				{Version: 4},
			}
			applied := []*migration.Migration{
				{Version: 1},
				{Version: 2},
			}

			filtered, err := Filter(available, applied)
			g.Assert(err).Equal(nil)
			g.Assert(len(filtered)).Equal(2)
			g.Assert(filtered[0].Version).Equal(uint64(3))
			g.Assert(filtered[1].Version).Equal(uint64(4))
		})

	})
}
