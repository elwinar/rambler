package configuration

import (
	. "github.com/franela/goblin"
	"github.com/spf13/pflag"
	"testing"
)

func TestGetEnvironment(t *testing.T) {
	g := Goblin(t)

	var s = "override"
	var i uint64 = 2
	var nilenv *Environment

	var override = RawEnvironment{
		Driver:    &s,
		Protocol:  &s,
		Host:      &s,
		Port:      &i,
		User:      &s,
		Password:  &s,
		Database:  &s,
		Directory: &s,
	}

	var empty = pflag.NewFlagSet("empty", pflag.ContinueOnError)
	empty.String("driver", "", "")
	empty.String("protocol", "", "")
	empty.String("host", "", "")
	empty.Uint("port", 0, "")
	empty.String("user", "", "")
	empty.String("password", "", "")
	empty.String("database", "", "")
	empty.String("directory", "", "")

	var cli = pflag.NewFlagSet("testing", pflag.ContinueOnError)
	cli.String("driver", "", "")
	cli.String("protocol", "", "")
	cli.String("host", "", "")
	cli.Uint("port", 0, "")
	cli.String("user", "", "")
	cli.String("password", "", "")
	cli.String("database", "", "")
	cli.String("directory", "", "")
	cli.Set("driver", "cli")
	cli.Set("protocol", "cli")
	cli.Set("host", "cli")
	cli.Set("port", "3")
	cli.Set("user", "cli")
	cli.Set("password", "cli")
	cli.Set("database", "cli")
	cli.Set("directory", "cli")

	var configuration = Configuration{
		Driver:    "default",
		Protocol:  "default",
		Host:      "default",
		Port:      1,
		User:      "default",
		Password:  "default",
		Database:  "default",
		Directory: "default",
		Environments: map[string]RawEnvironment{
			"default":  RawEnvironment{},
			"override": override,
		},
	}

	g.Describe("GetEnvironment", func() {
		g.It("Should reject empty environment", func() {
			env, err := GetEnvironment("", configuration, empty)
			g.Assert(err).Equal(ErrUnknownEnvironment)
			g.Assert(env).Equal(nilenv)
		})

		g.It("Should reject unknown environment", func() {
			env, err := GetEnvironment("error", configuration, empty)
			g.Assert(err).Equal(ErrUnknownEnvironment)
			g.Assert(env).Equal(nilenv)
		})

		g.It("Should use the defaults", func() {
			env, err := GetEnvironment("default", configuration, empty)
			g.Assert(err).Equal(nil)
			if env == nilenv {
				g.Fail("env equal <nil>")
			}
			g.Assert(env.Driver).Equal("default")
			g.Assert(env.Protocol).Equal("default")
			g.Assert(env.Host).Equal("default")
			g.Assert(env.Port).Equal(uint64(1))
			g.Assert(env.User).Equal("default")
			g.Assert(env.Password).Equal("default")
			g.Assert(env.Database).Equal("default")
			g.Assert(env.Directory).Equal("default")
		})

		g.It("Should override defaults with options", func() {
			env, err := GetEnvironment("override", configuration, empty)
			g.Assert(err).Equal(nil)
			if env == nilenv {
				g.Fail("env equal <nil>")
			}
			g.Assert(env.Driver).Equal("override")
			g.Assert(env.Protocol).Equal("override")
			g.Assert(env.Host).Equal("override")
			g.Assert(env.Port).Equal(uint64(2))
			g.Assert(env.User).Equal("override")
			g.Assert(env.Password).Equal("override")
			g.Assert(env.Database).Equal("override")
			g.Assert(env.Directory).Equal("override")
		})

		g.It("Should override options with cli", func() {
			env, err := GetEnvironment("override", configuration, cli)
			g.Assert(err).Equal(nil)
			if env == nilenv {
				g.Fail("env equal <nil>")
			}
			g.Assert(env.Driver).Equal("cli")
			g.Assert(env.Protocol).Equal("cli")
			g.Assert(env.Host).Equal("cli")
			g.Assert(env.Port).Equal(uint64(3))
			g.Assert(env.User).Equal("cli")
			g.Assert(env.Password).Equal("cli")
			g.Assert(env.Database).Equal("cli")
			g.Assert(env.Directory).Equal("cli")
		})
	})
}
