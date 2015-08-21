package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestReverse(t *testing.T) {
	var e = errors.New("error")
	var cases = []struct {
		initialized      bool
		initializedError error
		available        []*Migration
		availableError   error
		applied          []*Migration
		appliedError     error
		reverseError     error
		all              bool

		err      bool
		reversed []*Migration
	}{
		{
			initialized:      false,
			initializedError: nil,
			available:        nil,
			availableError:   nil,
			applied:          nil,
			appliedError:     nil,
			reverseError:     nil,
			all:              true,

			err:      true,
			reversed: nil,
		},
		{
			initialized:      false,
			initializedError: e,
			available:        nil,
			availableError:   nil,
			applied:          nil,
			appliedError:     nil,
			reverseError:     nil,
			all:              true,

			err:      true,
			reversed: nil,
		},
		{
			initialized:      true,
			initializedError: nil,
			available:        nil,
			availableError:   nil,
			applied:          nil,
			appliedError:     nil,
			reverseError:     nil,
			all:              true,

			err:      false,
			reversed: nil,
		},
		{
			initialized:      true,
			initializedError: nil,
			available: []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError:   nil,
			applied:          []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError:     nil,
			reverseError:     nil,
			all:              true,

			err:      false,
			reversed: []*Migration{
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
			},
		},
		{
			initialized:      true,
			initializedError: nil,
			available: []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError:   nil,
			applied:          []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError:     nil,
			reverseError:     nil,
			all:              false,

			err:      false,
			reversed: []*Migration{
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
			},
		},
		{
			initialized:      true,
			initializedError: nil,
			available: []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError:   nil,
			applied:          []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError:     nil,
			reverseError:     nil,
			all:              true,

			err:      true,
			reversed: nil,
		},
		{
			initialized:      true,
			initializedError: nil,
			available: []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "3.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError:   nil,
			applied:          []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "3.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError:     nil,
			reverseError:     nil,
			all:              true,

			err:      true,
			reversed: nil,
		},
		{
			initialized:      true,
			initializedError: nil,
			available: []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "3.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError:   nil,
			applied:          []*Migration{
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError:     nil,
			reverseError:     nil,
			all:              true,

			err:      false,
			reversed: []*Migration{
				&Migration{
					Name:   "2.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name:   "1.sql",
					reader: strings.NewReader(""),
				},
			},
		},
	}

	for n, c := range cases {
		var reversed []*Migration

		service := MockService{
			initialized: func() (bool, error) {
				return c.initialized, c.initializedError
			},
			available: func() ([]*Migration, error) {
				return c.available, c.availableError
			},
			applied: func() ([]*Migration, error) {
				return c.applied, c.appliedError
			},
			reverse: func(migration *Migration) error {
				reversed = append(reversed, migration)
				return c.reverseError
			},
		}

		err := reverse(service, c.all)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if !reflect.DeepEqual(reversed, c.reversed) {
			t.Error("case", n, "reversed the wrong migrations")
		}
	}
}
