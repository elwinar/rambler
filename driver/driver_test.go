package driver

import (
	"testing"
)

type MockDriver struct {
	new func(Config) (Conn, error)
}

func (d *MockDriver) New(c Config) (Conn, error) {
	return d.new(c)
}

func Test_Register_NilDriver(t *testing.T) {
	drivers = make(map[string]Driver)

	err := Register("test", nil)
	if err == nil {
		t.Fail()
	}
}

func Test_Register_AlreadyRegistered(t *testing.T) {
	drivers = make(map[string]Driver)

	err := Register("test", &MockDriver{})
	if err != nil {
		t.Error("unexpected error:", err)
	}

	err = Register("test", &MockDriver{})
	if err == nil {
		t.Fail()
	}
}

func Test_Register_OK(t *testing.T) {
	drivers = make(map[string]Driver)

	err := Register("test", &MockDriver{})
	if err != nil {
		t.Fail()
	}

	if _, found := drivers["test"]; !found {
		t.Fail()
	}
}

func Test_Get_NotRegistered(t *testing.T) {
	drivers = make(map[string]Driver)

	_, err := Get("test")
	if err == nil {
		t.Fail()
	}
}

func Test_Get_OK(t *testing.T) {
	drivers = make(map[string]Driver)
	err := Register("test", &MockDriver{})
	if err != nil {
		t.Fail()
	}

	_, err = Get("test")
	if err != nil {
		t.Fail()
	}
}
