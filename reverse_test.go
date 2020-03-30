package main

import (
	"errors"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/elwinar/rambler/log"
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
			all: true,
			err: true,
		},
		{
			initializedError: e,
			all:              true,
			err:              true,
		},
		{
			initialized: true,
			all:         true,
		},
		{
			initialized:    true,
			availableError: e,
			all:            true,
			err:            true,
		},
		{
			initialized:  true,
			appliedError: e,
			all:          true,
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
			all: true,
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
			all: true,
			err: true,
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
			all: true,
			err: true,
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
			all: true,
			reversed: []*Migration{
				{Name: "2.sql"},
				{Name: "1.sql"},
			},
		},
	}

	for n, c := range cases {
		var reversed []*Migration

		for i := range c.available {
			if c.available[i].reader == nil {
				c.available[i].reader = strings.NewReader("")
			}
		}

		for i := range c.applied {
			if c.applied[i].reader == nil {
				c.applied[i].reader = strings.NewReader("")
			}
		}

		for i := range c.reversed {
			if c.reversed[i].reader == nil {
				c.reversed[i].reader = strings.NewReader("")
			}
		}

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

		logger = log.NewLogger(func(l *log.Logger) {
			l.Output = ioutil.Discard
		})

		err := reverse(service, c.all, logger)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if !reflect.DeepEqual(reversed, c.reversed) {
			t.Error("case", n, "reversed the wrong migrations")
		}
	}
}
