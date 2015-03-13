package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	var e = errors.New("error")
	var cases = []struct{
		initialized      bool
		initializedError error
		initializeError  error
		available        []*Migration
		availableError   error
		applied          []*Migration
		appliedError     error
		applyError       error
		all              bool
		
		err              bool
		executed         []*Migration
	}{
		{
			initialized: true,
			initializedError: nil,
			initializeError: nil,
			available: nil,
			availableError: nil,
			applied: nil,
			appliedError: nil,
			applyError: nil,
			all: true,
			
			err: false,
			executed: nil,
		},
		{
			initialized: false,
			initializedError: e,
			initializeError: nil,
			available: nil,
			availableError: nil,
			applied: nil,
			appliedError: nil,
			applyError: nil,
			all: true,
			
			err: true,
			executed: nil,
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: e,
			available: nil,
			availableError: nil,
			applied: nil,
			appliedError: nil,
			applyError: nil,
			all: true,
			
			err: true,
			executed: nil,
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: nil,
			availableError: e,
			applied: nil,
			appliedError: nil,
			applyError: nil,
			all: true,
			
			err: true,
			executed: nil,
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: nil,
			availableError: nil,
			applied: nil,
			appliedError: e,
			applyError: nil,
			all: true,
			
			err: true,
			executed: nil,
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError: nil,
			applied: nil,
			appliedError: nil,
			applyError: e,
			all: true,
			
			err: true,
			executed: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
			},
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError: nil,
			applied: nil,
			appliedError: nil,
			applyError: nil,
			all: true,
			
			err: false,
			executed: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
			},
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError: nil,
			applied: nil,
			appliedError: nil,
			applyError: nil,
			all: false,
			
			err: false,
			executed: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
			},
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "zoo.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError: nil,
			applied: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError: nil,
			applyError: nil,
			all: false,
			
			err: false,
			executed: []*Migration{
				&Migration{
					Name: "zoo.sql",
					reader: strings.NewReader(""),
				},
			},
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "zoo.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError: nil,
			applied: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "zoo.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError: nil,
			applyError: nil,
			all: true,
			
			err: false,
			executed: nil,
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "zoo.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError: nil,
			applied: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "wee.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError: nil,
			applyError: nil,
			all: false,
			
			err: true,
			executed: nil,
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "zoo.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError: nil,
			applied: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "zoo.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError: nil,
			applyError: nil,
			all: false,
			
			err: true,
			executed: nil,
		},
		{
			initialized: false,
			initializedError: nil,
			initializeError: nil,
			available: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
			},
			availableError: nil,
			applied: []*Migration{
				&Migration{
					Name: "bar.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "foo.sql",
					reader: strings.NewReader(""),
				},
				&Migration{
					Name: "zoo.sql",
					reader: strings.NewReader(""),
				},
			},
			appliedError: nil,
			applyError: nil,
			all: false,
			
			err: true,
			executed: nil,
		},
	}
	
	for n, c := range cases {
		var executed []*Migration
		
		service := MockService{
			initialized: func()(bool, error) {
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
			apply: func(migration *Migration) error {
				executed = append(executed, migration)
				return c.applyError
			},
		}
		
		err := apply(service, c.all)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}
		
		if !reflect.DeepEqual(executed, c.executed) {
			t.Error("case", n, "executed the wrong migrations:", *executed[0])
		}
	}
}
