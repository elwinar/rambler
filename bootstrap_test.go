package main

import (
	"testing"
)

func TestBootstrap(t *testing.T) {
	var cases = []struct {
		configuration string
		environment   string
		err           bool
		initialized   bool
	}{
		{
			configuration: "test/valid.json",
			environment:   "default",
			err:           false,
			initialized:   true,
		},
		{
			configuration: "test/invalid.json",
			environment:   "default",
			err:           true,
			initialized:   false,
		},
		{
			configuration: "test/valid.json",
			environment:   "unknown",
			err:           true,
			initialized:   false,
		},
		{
			configuration: "test/faulty.json",
			environment:   "faulty",
			err:           true,
			initialized:   false,
		},
	}

	for n, c := range cases {
		service = nil

		err := bootstrap(c.configuration, c.environment)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if (service != nil) != c.initialized {
			t.Error("case", n, "got unexpected service:", service)
		}
	}
}
