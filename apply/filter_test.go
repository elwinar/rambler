package apply

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestFilter(t *testing.T) {
	g := Goblin(t)

	var missing []int64 = []int64{
		int64(1),
		int64(4),
	}

	var outOfOrder []int64 = []int64{
		int64(1),
		int64(2),
		int64(3),
		int64(4),
	}

	var applied []int64 = []int64{
		int64(1),
		int64(2),
		int64(4),
	}

	var available []int64 = []int64{
		int64(1),
		int64(2),
		int64(4),
		int64(5),
	}

	g.Describe("Filter", func() {
		g.It("Should complain about missing migrations", func() {
			filtered, err := Filter(missing, applied)
			g.Assert(err.Error()).Equal("missing migration 2")
			g.Assert(len(filtered)).Equal(0)
		})

		g.It("Should complain about out-of-order migrations", func() {
			filtered, err := Filter(outOfOrder, applied)
			g.Assert(err.Error()).Equal("out of order migration 3")
			g.Assert(len(filtered)).Equal(0)
		})

		g.It("Should return an empty slice when there is nothing to apply", func() {
			filtered, err := Filter(applied, applied)
			g.Assert(err).Equal(nil)
			g.Assert(len(filtered)).Equal(0)
		})

		g.It("Should return all remaining migrations where there is migrations to apply", func() {
			filtered, err := Filter(available, applied)
			g.Assert(err).Equal(nil)
			g.Assert(len(filtered)).Equal(1)
			g.Assert(filtered[0]).Equal(uint64(5))
		})

	})
}
