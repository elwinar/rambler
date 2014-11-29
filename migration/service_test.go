package migration

import (
	"errors"
	"github.com/elwinar/rambler/configuration"
	. "github.com/franela/goblin"
	"os"
	"testing"
)

var (
	nilService *service
)

func TestNewService(t *testing.T) {
	g := Goblin(t)
	
	mockEnv := configuration.Environment{
		Driver: "mock",
	}
	
	g.Describe("NewService", func() {
		g.It("Should reject unknown directory path", func() {
			s, err := newService(mockEnv, "", func(dir string) (os.FileInfo, error) {
				return nil, errors.New("error")
			}, func(env configuration.Environment) (Driver, error) {
				return nil, nil
			})
			g.Assert(err).Equal(ErrUnknownDirectory)
			g.Assert(s).Equal(nilService)
		})

		g.It("Should reject unknown driver", func() {
			s, err := newService(mockEnv, "", func(dir string) (os.FileInfo, error) {
				return nil, nil
			}, func(env configuration.Environment) (Driver, error) {
				return nil, errors.New("error")
			})
			g.Assert(err).Equal(ErrUnknownDriver)
			g.Assert(s).Equal(nilService)
		})

		g.It("Should return an initialized service", func() {
			d := &MockDriver{}
			s, err := newService(mockEnv, "dir", func(dir string) (os.FileInfo, error) {
				return nil, nil
			}, func(env configuration.Environment) (Driver, error) {
				return d, nil
			})
			g.Assert(err).Equal(nil)
			g.Assert(s.directory).Equal("dir")
			g.Assert(s.driver).Equal(d)
		})
	})
}
