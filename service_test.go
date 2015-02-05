package main

import (
	"testing"
)

func Test_NewService_InvalidDirectory(t *testing.T) {
	_, err := NewService(Environment{
		Driver: "mysql",
		Directory: "unknown",
	})

	if err == nil {
		t.Error("didn't failed on stat error")
	}
}

func Test_NewService_InvalidDriver(t *testing.T) {
	_, err := NewService(Environment{
		Driver: "unknown",
		Directory: "test",
	})

	if err == nil {
		t.Error("didn't failed on driver error")
	}
}

func Test_NewService_OK(t *testing.T) {
	s, err := NewService(Environment{
		Driver: "mysql",
		Directory: "test",
	})

	if err != nil {
		t.Error("unexpected error:", err)
	}

	if s == nil {
		t.Error("returned uninitialized service")
	}
}

func Test_Service_ListAvailableMigrations_ParseFilenames(t *testing.T) {
	
	s := CoreService{
		env: Environment{
			Directory: "test",
		},
	}
	versions := s.ListAvailableMigrations()
	
	if len(versions) != 8 {
		t.Errorf("didn't found the correct number of versions: %d", len(versions))
	}
}
