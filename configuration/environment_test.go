package configuration

import (
	. "github.com/franela/goblin"
	"github.com/spf13/pflag"
	"testing"
)

var (
	nilEnvironment *Environment
	override       string = "override"
	two            uint64 = 2

	overrideRawEnv RawEnvironment = RawEnvironment{
		Driver:     &override,
		Protocol:   &override,
		Host:       &override,
		Port:       &two,
		User:       &override,
		Password:   &override,
		Database:   &override,
		Migrations: &override,
	}

	emptyFlags *pflag.FlagSet
	cliFlags   *pflag.FlagSet
)

func init() {
	emptyFlags = pflag.NewFlagSet("empty", pflag.ContinueOnError)

	emptyFlags.String("driver", "", "")
	emptyFlags.String("protocol", "", "")
	emptyFlags.String("host", "", "")
	emptyFlags.Uint("port", 0, "")
	emptyFlags.String("user", "", "")
	emptyFlags.String("password", "", "")
	emptyFlags.String("database", "", "")
	emptyFlags.String("migrations", "", "")

	cliFlags = pflag.NewFlagSet("testing", pflag.ContinueOnError)

	cliFlags.String("driver", "", "")
	cliFlags.String("protocol", "", "")
	cliFlags.String("host", "", "")
	cliFlags.Uint("port", 0, "")
	cliFlags.String("user", "", "")
	cliFlags.String("password", "", "")
	cliFlags.String("database", "", "")
	cliFlags.String("migrations", "", "")

	cliFlags.Set("driver", "cli")
	cliFlags.Set("protocol", "cli")
	cliFlags.Set("host", "cli")
	cliFlags.Set("port", "3")
	cliFlags.Set("user", "cli")
	cliFlags.Set("password", "cli")
	cliFlags.Set("database", "cli")
	cliFlags.Set("migrations", "cli")
}

func TestGetEnvironment(t *testing.T) {
	g := Goblin(t)
	g.Describe("GetEnvironment", func() {
		g.It("Should reject empty environment", func() {
			env, err := GetEnvironment("", Configuration{}, emptyFlags)
			g.Assert(env).Equal(nilEnvironment)
			g.Assert(err).Equal(ErrUnknownEnvironment)
		})

		g.It("Should reject unknown environment", func() {
			env, err := GetEnvironment("error", Configuration{}, emptyFlags)
			g.Assert(env).Equal(nilEnvironment)
			g.Assert(err).Equal(ErrUnknownEnvironment)
		})

		g.It("Should use the defaults", func() {
			env, err := GetEnvironment("override", Configuration{
				Driver:     "default",
				Protocol:   "default",
				Host:       "default",
				Port:       1,
				User:       "default",
				Password:   "default",
				Database:   "default",
				Migrations: "default",
				Environments: map[string]RawEnvironment{
					"override": RawEnvironment{},
				},
			}, emptyFlags)
			if env == nil {
				g.Fail("env equal <nil>")
			}
			g.Assert(env.Driver).Equal("default")
			g.Assert(env.Protocol).Equal("default")
			g.Assert(env.Host).Equal("default")
			g.Assert(env.Port).Equal(uint64(1))
			g.Assert(env.User).Equal("default")
			g.Assert(env.Password).Equal("default")
			g.Assert(env.Database).Equal("default")
			g.Assert(env.Migrations).Equal("default")
			g.Assert(err).Equal(nil)
		})

		g.It("Should override defaults with options", func() {
			env, err := GetEnvironment("override", Configuration{
				Driver:     "default",
				Protocol:   "default",
				Host:       "default",
				Port:       1,
				User:       "default",
				Password:   "default",
				Database:   "default",
				Migrations: "default",
				Environments: map[string]RawEnvironment{
					"override": overrideRawEnv,
				},
			}, emptyFlags)
			if env == nil {
				g.Fail("env equal <nil>")
			}
			g.Assert(env.Driver).Equal("override")
			g.Assert(env.Protocol).Equal("override")
			g.Assert(env.Host).Equal("override")
			g.Assert(env.Port).Equal(uint64(2))
			g.Assert(env.User).Equal("override")
			g.Assert(env.Password).Equal("override")
			g.Assert(env.Database).Equal("override")
			g.Assert(env.Migrations).Equal("override")
			g.Assert(err).Equal(nil)
		})

		g.It("Should override options with cli", func() {
			env, err := GetEnvironment("override", Configuration{
				Driver:     "default",
				Protocol:   "default",
				Host:       "default",
				Port:       1,
				User:       "default",
				Password:   "default",
				Database:   "default",
				Migrations: "default",
				Environments: map[string]RawEnvironment{
					"override": overrideRawEnv,
				},
			}, cliFlags)
			if env == nil {
				g.Fail("env equal <nil>")
			}
			g.Assert(env.Driver).Equal("cli")
			g.Assert(env.Protocol).Equal("cli")
			g.Assert(env.Host).Equal("cli")
			g.Assert(env.Port).Equal(uint64(3))
			g.Assert(env.User).Equal("cli")
			g.Assert(env.Password).Equal("cli")
			g.Assert(env.Database).Equal("cli")
			g.Assert(env.Migrations).Equal("cli")
			g.Assert(err).Equal(nil)
		})
	})
}
