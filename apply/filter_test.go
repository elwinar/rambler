package apply

import (
	"errors"
	"testing"
)

func TestFilter(t *testing.T) {
	type Case struct {
		Available []uint64
		Applied   []uint64
		Filtered  []uint64
		Error     error
	}
	var cases = []Case{
		Case{
			Available: []uint64{1, 4},
			Applied:   []uint64{1, 2, 4},
			Filtered:  []uint64{},
			Error:     errors.New("missing migration 2"),
		},
		Case{
			Available: []uint64{1, 2},
			Applied:   []uint64{1, 2, 4},
			Filtered:  []uint64{},
			Error:     errors.New("missing migration 4"),
		},
		Case{
			Available: []uint64{1, 2, 3, 4},
			Applied:   []uint64{1, 2, 4},
			Filtered:  []uint64{},
			Error:     errors.New("out of order migration 3"),
		},
		Case{
			Available: []uint64{1, 4},
			Applied:   []uint64{1, 2, 4},
			Filtered:  []uint64{},
			Error:     errors.New("missing migration 2"),
		},
		Case{
			Available: []uint64{1, 2, 4, 5},
			Applied:   []uint64{1, 2, 4},
			Filtered:  []uint64{5},
			Error:     nil,
		},
	}

	for n, c := range cases {
		filtered, err := Filter(c.Available, c.Applied)

		if err == nil && c.Error != nil {
			t.Errorf("expected error on case %d", n)
		} else if err != nil && c.Error == nil {
			t.Errorf("unexpected error on case %d: %s", n, err.Error())
		} else if err != nil && c.Error != nil && err.Error() != c.Error.Error() {
			t.Errorf("didn't returned expected error on case %d: %s", n, err.Error())
		}

		if len(filtered) != len(c.Filtered) {
			t.Error("incorrectly filtered case %d: kept %d", n, len(filtered))
			continue
		}

		for i := 0; i < len(filtered); i++ {
			if filtered[i] != c.Filtered[i] {
				t.Error("incorrect migration version on index %d of case %d: %d", i, n, filtered[i])
			}
		}
	}
}
