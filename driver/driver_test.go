package driver

import (
	"errors"
	"testing"
)

type MockDriver struct {
	new func(string, string) (Conn, error)
}

func (d *MockDriver) New(dsn, schema string) (Conn, error) {
	return d.new(dsn, schema)
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
	
	_, err := Get("test", "", "")
	if err == nil {
		t.Fail()
	}
}

func Test_Get_InitializeError(t *testing.T) {
	drivers = make(map[string]Driver)
	
	driver := &MockDriver{}
	driver.new = func(_, _ string) (Conn, error) {
		return nil, errors.New("initialize error")
	}

	err := Register("test", driver)
	if err != nil {
		t.Fail()
	}

	_, err = Get("test", "", "")
	if err == nil {
		t.Fail()
	}
}

func Test_Get_OK(t *testing.T) {
	drivers = make(map[string]Driver)
	driver := &MockDriver{}
	driver.new = func(_, _ string) (Conn, error) {
		return nil, nil
	}

	err := Register("test", driver)
	if err != nil {
		t.Fail()
	}

	_, err = Get("test", "", "")
	if err != nil {
		t.Fail()
	}
}
