package main

import (
	"testing"
)

func Test_Environment_DSN_Unknown(t *testing.T) {
	e := Environment{
		Driver: "unknown",
	}
	
	if e.DSN() != "" {
		t.Fail()
	}
}
