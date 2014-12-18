package driver

import (
	"errors"
	"github.com/elwinar/rambler/configuration"
	"testing"
)

func Test_Register_NilDriver(t *testing.T) {
	drivers := make(map[string]Driver)

	err := register("test", nil, drivers)
	if err == nil {
		t.Error("didn't failed on nil driver")
		return
	}

	if err.Error() != "not a valid driver" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_Register_AlreadyRegistered(t *testing.T) {
	drivers := make(map[string]Driver)

	err := register("test", &MockDriver{}, drivers)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	err = register("test", &MockDriver{}, drivers)
	if err == nil {
		t.Error("didn't failed on already registered driver")
		return
	}

	if err.Error() != "driver test already registered" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_Register_OK(t *testing.T) {
	drivers := make(map[string]Driver)

	err := register("test", &MockDriver{}, drivers)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if _, found := drivers["test"]; !found {
		t.Error("didn't registered driver")
	}
}

func Test_Get_NotRegistered(t *testing.T) {
	drivers := make(map[string]Driver)
	env := configuration.Environment{
		Driver: "test",
	}

	_, err := get(env, drivers)
	if err == nil {
		t.Error("didn't failed on unregistered driver")
		return
	}

	if err.Error() != "driver test not registered" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_Get_InitializeError(t *testing.T) {
	drivers := make(map[string]Driver)
	driver := &MockDriver{}
	driver.new = func(env configuration.Environment) (Conn, error) {
		return nil, errors.New("initialize error")
	}
	
	env := configuration.Environment{
		Driver: "test",
	}

	err := register("test", driver, drivers)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	_, err = get(env, drivers)
	if err == nil {
		t.Error("didn't failed on initialize error")
	}

	if err.Error() != "initialize error" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_Get_OK(t *testing.T) {
	drivers := make(map[string]Driver)
	conn := &MockConn{}
	driver := &MockDriver{}
	driver.new = func(env configuration.Environment) (Conn, error) {
		return conn, nil
	}
	
	env := configuration.Environment{
		Driver: "test",
	}

	err := register("test", driver, drivers)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	c, err := get(env, drivers)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if c != conn {
		t.Error("didn't returned expected driver")
	}
}
