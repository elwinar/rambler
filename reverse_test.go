package main

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/elwinar/rambler/log"
)

func TestReverse(t *testing.T) {
	e := errors.New("error")
	cases := []struct {
		initialized      bool
		initializedError error
		available        []*Migration
		availableError   error
		applied          []*Migration
		appliedError     error
		reverseError     error
		all              bool
		save             bool
		migration        string

		err      bool
		reversed []*Migration
	}{
		{
			all:  true,
			save: true,
			err:  true,
		},
		{
			initializedError: e,
			all:              true,
			save:             true,
			err:              true,
		},
		{
			initialized: true,
			all:         true,
			save:        true,
		},
		{
			initialized:    true,
			availableError: e,
			all:            true,
			save:           true,
			err:            true,
		},
		{
			initialized:  true,
			appliedError: e,
			all:          true,
			save:         true,
			err:          true,
		},
		{
			initialized: true,
			available: []*Migration{
				{Name: "1.sql"},
				{Name: "2.sql"},
			},
			availableError: nil,
			applied: []*Migration{
				{Name: "1.sql"},
				{Name: "2.sql"},
			},
			all:  true,
			save: true,
			reversed: []*Migration{
				{Name: "2.sql"},
				{Name: "1.sql"},
			},
		},
		{
			initialized: true,
			available: []*Migration{
				{Name: "1.sql"},
				{Name: "2.sql"},
			},
			applied: []*Migration{
				{Name: "1.sql"},
				{Name: "2.sql"},
			},
			save: true,
			reversed: []*Migration{
				{Name: "2.sql"},
			},
		},
		{
			initialized: true,
			available: []*Migration{
				{Name: "1.sql"},
			},
			applied: []*Migration{
				{Name: "1.sql"},
				{Name: "2.sql"},
			},
			all:  true,
			save: true,
			err:  true,
		},
		{
			initialized: true,
			available: []*Migration{
				{Name: "1.sql"},
				{Name: "2.sql"},
				{Name: "3.sql"},
			},
			applied: []*Migration{
				{Name: "1.sql"},
				{Name: "3.sql"},
			},
			all:  true,
			save: true,
			err:  true,
		},
		{
			initialized: true,
			available: []*Migration{
				{Name: "1.sql"},
				{Name: "2.sql"},
				{Name: "3.sql"},
			},
			availableError: nil,
			applied: []*Migration{
				{Name: "1.sql"},
				{Name: "2.sql"},
			},
			all:  true,
			save: true,
			reversed: []*Migration{
				{Name: "2.sql"},
				{Name: "1.sql"},
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
			reverse: func(migration *Migration, save bool) error {
				if save {
					reversed = append(reversed, migration)
				}
				return c.reverseError
			},
		}

		logger = log.NewLogger(func(l *log.Logger) {
			l.Output = io.Discard
		})

		err := reverse(service, c.all, c.save, c.migration, logger)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if !reflect.DeepEqual(reversed, c.reversed) {
			t.Error("case", n, "reversed the wrong migrations")
		}
	}
}
