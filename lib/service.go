package lib

import (
	drv "github.com/elwinar/rambler/lib/driver"
	_ "github.com/elwinar/rambler/lib/driver/mysql"
	"os"
)

type Service interface {
	drv.Driver
}

type service struct {
	driver    drv.Driver
	directory string
}

func NewService(driver, dsn, directory string) (Service, error) {
	return newService(driver, dsn, directory, os.Stat, drv.Get)
}

type stater func(string) (os.FileInfo, error)
type driverConstructor func(string, string) (drv.Driver, error)

func newService(driverName, dsn, directory string, stat stater, newDriver driverConstructor) (*service, error) {
	if _, err := stat(directory); err != nil {
		return nil, ErrUnknownDirectory
	}

	driver, err := newDriver(driverName, dsn)
	if err != nil {
		return nil, ErrUnknownDriver
	}

	return &service{
		driver:    driver,
		directory: directory,
	}, nil
}

func (s service) MigrationTableExists() (bool, error) {
	return s.driver.MigrationTableExists()
}

func (s service) CreateMigrationTable() error {
	return s.driver.CreateMigrationTable()
}
