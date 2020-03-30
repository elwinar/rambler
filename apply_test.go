package main

import (
	"errors"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/elwinar/rambler/log"
)

func TestApply(t *testing.T) {
	var e = errors.New("error")
	var cases = []struct {
		initialized      bool
		initializedError error
		initializeError  error
		available        []*Migration
		availableError   error
		applied          []*Migration
		appliedError     error
		applyError       error
		all              bool
		save             bool
		migration        string

		err      bool
		executed []*Migration
	}{
		{
			initialized: true,
			all:         true,
			save:        true,
		},
		{
			initializedError: e,
			all:              true,
			save:             true,
			err:              true,
		},
		{
			initializeError: e,
			all:             true,
			save:            true,
			err:             true,
		},
		{
			availableError: e,
			all:            true,
			save:           true,
			err:            true,
		},
		{
			appliedError: e,
			all:          true,
			save:         true,
			err:          true,
		},
		{
			available: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
			},
			applyError: e,
			all:        true,
			save:       true,
			err:        true,
			executed: []*Migration{
				{Name: "bar.sql"},
			},
		},
		{
			available: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
			},
			all:  true,
			save: true,
			executed: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
			},
		},
		{
			available: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
			},
			save: true,
			executed: []*Migration{
				{Name: "bar.sql"},
			},
		},
		{
			available: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
				{Name: "zoo.sql"},
			},
			applied: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
			},
			save: true,
			executed: []*Migration{
				{Name: "zoo.sql"},
			},
		},
		{
			available: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
				{Name: "zoo.sql"},
			},
			applied: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
				{Name: "zoo.sql"},
			},
			all:  true,
			save: true,
		},
		{
			available: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
				{Name: "zoo.sql"},
			},
			applied: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
				{Name: "wee.sql"},
			},
			save: true,
			err:  true,
		},
		{
			available: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
				{Name: "zoo.sql"},
			},
			applied: []*Migration{
				{Name: "bar.sql"},
				{Name: "zoo.sql"},
			},
			save: true,
			err:  true,
		},
		{
			available: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
			},
			applied: []*Migration{
				{Name: "bar.sql"},
				{Name: "foo.sql"},
				{Name: "zoo.sql"},
			},
			save: true,
			err:  true,
		},
	}

	for n, c := range cases {
		var executed []*Migration

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

		for i := range c.executed {
			if c.executed[i].reader == nil {
				c.executed[i].reader = strings.NewReader("")
			}
		}

		service := MockService{
			initialized: func() (bool, error) {
				return c.initialized, c.initializedError
			},
			initialize: func() error {
				return c.initializeError
			},
			available: func() ([]*Migration, error) {
				return c.available, c.availableError
			},
			applied: func() ([]*Migration, error) {
				return c.applied, c.appliedError
			},
			apply: func(migration *Migration, save bool) error {
				if save {
					executed = append(executed, migration)
				}
				return c.applyError
			},
		}

		logger = log.NewLogger(func(l *log.Logger) {
			l.Output = ioutil.Discard
		})

		err := apply(service, c.all, c.save, c.migration, logger)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if !reflect.DeepEqual(executed, c.executed) {
			t.Errorf("#%d: executed the wrong migrations: wanted %+v, got %+v", n, c.executed, executed)
		}
	}
}
