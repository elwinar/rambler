package main

import (
	"errors"
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/driver"
	"os"
	"testing"
)

func Test_NewService_InvalidDirectory(t *testing.T) {
	_, err := newService(configuration.Environment{
		Directory: "test",
	}, func(_ string) (os.FileInfo, error) {
		return nil, errors.New("stat error")
	}, func(_ configuration.Environment) (driver.Conn, error) {
		return &MockConn{}, nil
	})

	if err == nil {
		t.Error("didn't failed on stat error")
	}

	if err.Error() != "directory test unavailable: stat error" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_NewService_InvalidDriver(t *testing.T) {
	_, err := newService(configuration.Environment{
		Driver: "test",
	}, func(_ string) (os.FileInfo, error) {
		return nil, nil
	}, func(_ configuration.Environment) (driver.Conn, error) {
		return &MockConn{}, errors.New("driver error")
	})

	if err == nil {
		t.Error("didn't failed on driver error")
	}

	if err.Error() != "unable to initialize driver test: driver error" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_NewService_OK(t *testing.T) {
	s, err := newService(configuration.Environment{
		Driver: "test",
	}, func(_ string) (os.FileInfo, error) {
		return nil, nil
	}, func(_ configuration.Environment) (driver.Conn, error) {
		return &MockConn{}, nil
	})

	if err != nil {
		t.Error("unexpected error:", err)
	}

	if s == nil {
		t.Error("returned uninitialized service")
	}
}

func Test_Service_ListAvailableMigrations_GlobError(t *testing.T) {
	_, err := listAvailableMigrations(configuration.Environment{
		Directory: "test",
	}, func(_ string) ([]string, error) {
		return nil, errors.New("glob error")
	})

	if err == nil {
		t.Error("didn't failed on driver error")
	}

	if err.Error() != "directory test unavailable: glob error" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_Service_ListAvailableMigrations_ParseFilenames(t *testing.T) {
	type Case struct {
		Files    []string
		Versions []uint64
	}
	var cases = []Case{
		Case{
			Files:    []string{"0_one.sql", "1_two.sql"},
			Versions: []uint64{0, 1},
		},
		Case{
			Files:    []string{"one_one.sql", "1_two.sql"},
			Versions: []uint64{1},
		},
		Case{
			Files:    []string{"0.5_one.sql", "1_two.sql"},
			Versions: []uint64{1},
		},
		Case{
			Files:    []string{"one.sql", "1_two.sql"},
			Versions: []uint64{1},
		},
	}

	for n, c := range cases {
		versions, err := listAvailableMigrations(configuration.Environment{
			Directory: "test",
		}, func(_ string) ([]string, error) {
			return c.Files, nil
		})

		if err != nil {
			t.Error("unexpected error:", err)
		}

		if len(versions) != len(c.Versions) {
			t.Log(versions)
			t.Errorf("didn't found the correct number of versions for case %d: %d != %d", n, len(versions), len(c.Versions))
			continue
		}

		for i := 0; i < len(versions); i++ {
			if versions[i] != c.Versions[i] {
				t.Error("didn't parsed correctly file %d of case %d", i, n)
			}
		}
	}
}
